// consider EOF and Semicolon equivalent - Terminator or something
package dsql

import (
	"io"
	"regexp"
	"strings"
	"text/scanner"
)

func NewLexer(source io.Reader) *Lexer {
	l := &Lexer{}
	l.pos = -1
	l.scn.Init(source)
	l.scan()
	return l
}

type Lexer struct {
	scn     scanner.Scanner
	pos     int
	tokens  []Token
	strings []string
}

func (l *Lexer) scan() {
	var t Token
	for t != EOF {
		t = l.tokenize(l.sNext())
		l.tokens = append(l.tokens, t)
		l.strings = append(l.strings, l.sText())
	}
}

func (l *Lexer) Peek() Token {
	if l.pos < 0 {
		l.Next()
	}
	return l.tokens[l.pos]
}

func (l *Lexer) Next() Token {
	l.pos++
	return l.tokens[l.pos]
}

func (l *Lexer) Text() string {
	return l.strings[l.pos]
}

func (l *Lexer) tokenize(r rune) (t Token) {
	switch r {
	case scanner.Ident:
		if l.match(Keywords) {
			t = Keyword
		} else if l.match(Types) {
			t = Type
		} else if l.match(Constraints) {
			t = Constraint
		} else if l.match(Operators) {
			t = Operator
		} else {
			t = Identifier
		}
	case scanner.Float, scanner.Int:
		t = Number
	case scanner.String:
		t = String
	case scanner.EOF:
		t = EOF
	default:
		if l.match(Operators) {
			t = Operator
		} else {
			switch l.sText() {
			case "*":
				t = Wildcard
			case ",":
				t = Comma
			case ";":
				t = EOF
			case "(":
				t = LeftParen
			case ")":
				t = RightParen
			default:
				t = Unknown
			}
		}
	}
	return t
}

func (l *Lexer) sNext() rune {
	return l.scn.Scan()
}

func (l *Lexer) sText() string {
	return strings.ToLower(l.scn.TokenText())
}

func (l *Lexer) match(s string) bool {
	matched, err := regexp.MatchString(s, l.sText())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) is(s string) bool {
	return l.sText() == s
}
