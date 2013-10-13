package dsql

import (
	"errors"
	"log"

	"io"

	"regexp"
	"strings"
	"text/scanner"
)

type token rune

const (
	TokKeyword token = iota
	TokId
	TokWildcard
	TokString
	TokNumber
	TokOp
	TokComma
	TokSemicolon
	TokEOF
	TokUnknown
)

var ErrUnexpectedToken = errors.New("lexer: unexpected token")

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
		names = append(names, l.Name(t))
	}
	return names
}

func (l *Lexer) Name(t token) string {
	var name string
	switch t {
	case TokKeyword:
		name = "Keyword"
	case TokId:
		name = "Id"
	case TokWildcard:
		name = "Wildcard"
	case TokString:
		name = "String"
	case TokNumber:
		name = "Number"
	case TokOp:
		name = "Op"
	case TokComma:
		name = "Comma"
	case TokSemicolon:
		name = "Semicolon"
	case TokEOF:
		name = "EOF"
	case TokUnknown:
		name = "Unknown"
	}
	return name
}

func (l *Lexer) Tokens() (tokens []token) {
	var t token
	for {
		t = l.Next()
		tokens = append(tokens, t)
		if t == TokEOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) Next() (t token) {
	switch l.scn.Scan() {
	case scanner.Ident:
		if l.isKeyword() {
			t = TokKeyword
		} else {
			t = TokId
		}
	case scanner.Float, scanner.Int:
		t = TokNumber
	case scanner.String:
		t = TokString
	case scanner.EOF:
		t = TokEOF
	default:
		if l.isWildcard() {
			t = TokWildcard
		} else if l.isComma() {
			t = TokComma
		} else if l.isOp() {
			t = TokOp
		} else {
			log.Print("unknown: ", l.NextString())
			t = TokUnknown
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

func (l *Lexer) isOp() bool {
	matched, err := regexp.MatchString("=|!=|>|<", l.NextString())
	if err != nil {
		return false
	}
	return matched
}

func (l *Lexer) scanString(raw rune) token {
	return TokEOF
}

func (l *Lexer) NextString() string {
	return strings.ToLower(l.scn.TokenText())
}
