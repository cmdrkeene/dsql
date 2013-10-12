package dsql

import (
	"errors"
	"strings"

	"io"

	"text/scanner"
)

type token int

const (
	TokCreate token = iota
	TokSelect
	TokInsert
	TokUpdate
	TokDelete
	TokFrom
	TokWhere
	TokEqual
	TokNotEqual // !=, <>
	TokEOF
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

func (l *Lexer) Next() (t token, err error) {
	l.scn.Scan()
	switch strings.ToLower(l.scn.TokenText()) {
	case "select":
		t = TokSelect
	case "from":
		t = TokFrom
	case "where":
		t = TokWhere
	default:
		err = ErrUnexpectedToken
	}

	if err != nil {
		return 0, err
	}

	return t, nil
}
