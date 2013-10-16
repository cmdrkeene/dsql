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
		if l.isKeyword() {
			t = Keyword
		} else if l.isType() {
			t = Type
		} else if l.isOperator() {
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
		if l.isWildcard() {
			t = Wildcard
		} else if l.isComma() {
			t = Comma
		} else if l.isSemicolon() {
			t = Semicolon
		} else if l.isOperator() {
			t = Operator
		} else if l.isLeftParen() {
			t = LeftParen
		} else if l.isRightParen() {
			t = RightParen
		} else {
			t = Unknown
		}
	}
	return t
}

func (l *Lexer) isKeyword() bool {
	matched, err := regexp.MatchString(Keywords, l.sText())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) isType() bool {
	matched, err := regexp.MatchString(Types, l.sText())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) isId() bool {
	return true
}

func (l *Lexer) isWildcard() bool {
	return l.sText() == "*"
}

func (l *Lexer) isComma() bool {
	return l.sText() == ","
}

func (l *Lexer) isSemicolon() bool {
	return l.sText() == ";"
}

func (l *Lexer) isLeftParen() bool {
	return l.sText() == "("
}

func (l *Lexer) isRightParen() bool {
	return l.sText() == ")"
}

func (l *Lexer) isOperator() bool {
	matched, err := regexp.MatchString(Operators, l.sText())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) sNext() rune {
	return l.scn.Scan()
}

func (l *Lexer) sText() string {
	return strings.ToLower(l.scn.TokenText())
}
