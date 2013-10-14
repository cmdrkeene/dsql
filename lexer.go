package dsql

import (
	"io"

	"regexp"
	"strings"
	"text/scanner"
)

type token rune

const (
	Keyword token = iota
	Identifier
	Wildcard
	String
	Number
	Operator
	Comma
	Semicolon
	EOF
	Unknown
)

// shorthand
const (
	K = Keyword
	I = Identifier
	W = Wildcard
	N = Number
	S = String
	O = Operator
	C = Comma
	M = Semicolon
	E = EOF
	U = Unknown
)

var TokenNames = map[token]string{
	Keyword:    "Keyword",
	Identifier: "Identifier",
	Wildcard:   "Wildcard",
	String:     "String",
	Number:     "Number",
	Operator:   "Operator",
	Comma:      "Comma",
	Semicolon:  "Semicolon",
	EOF:        "EOF",
	Unknown:    "Unkown",
}

func NewLexer(source io.Reader) *Lexer {
	l := &Lexer{}
	l.scn.Init(source)
	return l
}

type Lexer struct {
	scn scanner.Scanner
	tok rune
}

func (l *Lexer) Names(tokens []token) (names []string) {
	for _, t := range tokens {
		names = append(names, TokenNames[t])
	}
	return names
}

func (l *Lexer) Tokens() (tokens []token) {
	var t token
	for {
		t = l.Next()
		tokens = append(tokens, t)
		if t == EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) Next() (t token) {
	switch l.scn.Scan() {
	case scanner.Ident:
		if l.isKeyword() {
			t = Keyword
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
		} else if l.isOperator() {
			t = Operator
		} else {
			t = Unknown
		}
	}
	return t
}

func (l *Lexer) isKeyword() bool {
	matched, err := regexp.MatchString("select|insert|create|update|delete|from|where|and", l.NextString())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) isId() bool {
	return true
}

func (l *Lexer) isWildcard() bool {
	return l.NextString() == "*"
}

func (l *Lexer) isComma() bool {
	return l.NextString() == ","
}

func (l *Lexer) isOperator() bool {
	matched, err := regexp.MatchString("=|!=|>|<", l.NextString())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) NextString() string {
	return strings.ToLower(l.scn.TokenText())
}
