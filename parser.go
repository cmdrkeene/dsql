package dsql

import (
	"errors"
	"fmt"

	"strings"
)

var ErrUnexpectedToken = errors.New("parser: unexpected token")

func Parse(source string) (interface{}, error) {
	parser := &Parser{
		src: source,
		lex: NewLexer(strings.NewReader(source)),
	}
	return parser.Parse()
}

type Parser struct {
	src string
	lex *Lexer
	err error
}

func (p *Parser) Parse() (interface{}, error) {
	var req interface{}

	switch p.next() {
	case Keyword:
		switch p.text() {
		case "select":
			req = p.Select()
		case "insert":
			req = p.Insert()
		}
	}
	return req, p.err
}

func (p *Parser) Select() interface{} {
	req := Query{}

	switch p.next() {
	case Wildcard:
		req.Select = "ALL_ATTRIBUTES"
	case Identifier:
		var ids []string
		var t Token
		for t != Keyword {
			t = p.peek()
			if t == Identifier {
				ids = append(ids, p.text())
				p.next()
			} else if t == Comma {
				p.next()
			} else {
				p.errUnexpected()
				return nil
			}
		}
	default:
		p.errUnexpected()
		return nil
	}

	if p.next() != Keyword && p.text() != "from" {
		p.errUnexpected()
		return nil
	}

	if p.next() != Identifier {
		p.errUnexpected()
		return nil
	}

	req.TableName = p.text()

	return req
}

func (p *Parser) Insert() interface{} {
	return nil
}

func (p *Parser) next() Token {
	return p.lex.Next()
}

func (p *Parser) peek() Token {
	return p.lex.Peek()
}

func (p *Parser) text() string {
	return p.lex.Text()
}

func (p *Parser) errUnexpected() {
	p.err = fmt.Errorf("parser: unexpected token '%s' in '%s', expected %s", p.text(), p.src, Names[p.peek()])
}
