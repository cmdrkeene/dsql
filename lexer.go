package dsql

import (
	"errors"
	"fmt"
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

func (l *Lexer) String(t token) string {
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
	return fmt.Sprintf("%s (%s)", name, l.NextString())
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
	raw := l.scn.Scan()

	switch raw {
	case scanner.Ident:
		if l.isKeyword() {
			t = TokKeyword
		} else {
			t = TokId
		}
	case scanner.EOF:
		t = TokEOF
	default:
		if l.isWildcard() {
			t = TokWildcard
		} else {
			t = TokUnknown
		}
	}
	return t
}

func (l *Lexer) isKeyword() bool {
	matched, err := regexp.MatchString("select|insert|create|update|delete|from|where", l.NextString())
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

func (l *Lexer) scanString(raw rune) token {
	return TokEOF
}

func (l *Lexer) NextString() string {
	return strings.ToLower(l.scn.TokenText())
}
