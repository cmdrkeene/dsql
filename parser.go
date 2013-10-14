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

// get wildcard or identifier
// if comma, add identifiers until next keyword
// if keyword is not from, error
// if next is not identifier, error
func (p *Parser) Select() interface{} {
	req := Query{}

	p.next()

	if !p.match(Wildcard, Identifier) {
		return nil
	}

	if p.peek() == Wildcard {
		p.next()
	}

	if p.peek() == Identifier {
		ids := []string{p.text()}
		for p.peek() == Identifier || p.peek() == Comma {
			if p.next() == Identifier {
				ids = append(ids, p.text())
			}
		}
		req.AttributesToGet = ids
	}

	if !p.match(Keyword) {
		return nil
	}

	if !p.matchText("from") {
		return nil
	}

	p.next()
	if !p.match(Identifier) {
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

func (p *Parser) match(tokens ...Token) bool {
	for _, t := range tokens {
		if p.peek() == t {
			return true
		}
	}
	p.errUnexpected()
	return false
}

func (p *Parser) matchText(t string) bool {
	if p.text() == t {
		return true
	}
	p.errUnexpected()
	return false
}

func (p *Parser) errUnexpected() {
	p.err = fmt.Errorf("parser: unexpected token '%s' in '%s' expected %s", p.text(), p.src, Names[p.peek()])
}
