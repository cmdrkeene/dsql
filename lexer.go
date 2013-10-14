package dsql

import (
	"io"
	"regexp"
	"strings"
	"text/scanner"
)

func NewLexer(source io.Reader) *Lexer {
	l := &Lexer{}
	l.scn.Init(source)
	return l
}

type Lexer struct {
	scn scanner.Scanner
	tok rune
}

func (l *Lexer) Tokens() (tokens []Token) {
	var t Token
	for {
		t = l.Next()
		tokens = append(tokens, t)
		if t == EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) Peek() (t Token) {
	return l.tokenize(l.scn.Peek())
}

func (l *Lexer) Next() (t Token) {
	return l.tokenize(l.scn.Scan())
}

func (l *Lexer) tokenize(r rune) (t Token) {
	switch r {
	case scanner.Ident:
		if l.isKeyword() {
			t = Keyword
		} else if l.isType() {
			t = Type
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
	matched, err := regexp.MatchString(Keywords, l.Text())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) isType() bool {
	matched, err := regexp.MatchString(Types, l.Text())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) isId() bool {
	return true
}

func (l *Lexer) isWildcard() bool {
	return l.Text() == "*"
}

func (l *Lexer) isComma() bool {
	return l.Text() == ","
}

func (l *Lexer) isSemicolon() bool {
	return l.Text() == ";"
}

func (l *Lexer) isLeftParen() bool {
	return l.Text() == "("
}

func (l *Lexer) isRightParen() bool {
	return l.Text() == ")"
}

func (l *Lexer) isOperator() bool {
	matched, err := regexp.MatchString(Operators, l.Text())
	if err != nil {
		return false
	}
	return matched
}

// returns lowercased string data in scanners current token
func (l *Lexer) Text() string {
	return strings.ToLower(l.scn.TokenText())
}
