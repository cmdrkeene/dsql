package dsql

import (
	"errors"
	"fmt"
	"log"

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

// recover from unexpected eof
// recover from unexpected token
func (p *Parser) Parse() (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	var req interface{}

	switch p.token() {
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
	query := Query{}

	p.matchS(Keyword, "select")

	if p.token() == Wildcard {
		p.consume()
	} else {
		query.AddColumn(p.match(Identifier))

		for p.token() == Comma {
			p.match(Comma)
			id := p.text()
			p.match(Identifier)
			query.AddColumn(id)
		}
	}

	p.matchS(Keyword, "from")

	query.TableName = p.match(Identifier)

	if p.token() == EOF {
		return query
	}

	// WHERE (Ident Operator (Number|String))+ (EOF|Keyword)

	p.matchS(Keyword, "where")

	query.AddCondition(p.expr())

	for p.token() == Keyword && p.text() == "and" {
		p.matchS(Keyword, "and")
		query.AddCondition(p.expr())
	}

	return query
}

type Expression struct {
	Identifier string
	Operator   string
	ValueToken Token
	ValueText  string
}

func (p *Parser) expr() Expression {
	return Expression{
		p.match(Identifier),
		p.match(Operator),
		p.token(),
		p.trim(p.match(String, Number)),
	}
}

func (p *Parser) Insert() interface{} {
	return nil
}

func (p *Parser) consume() Token {
	return p.lex.Next()
}

func (p *Parser) token() Token {
	return p.lex.Peek()
}

func (p *Parser) text() string {
	return p.lex.Text()
}

func (p *Parser) match(tokens ...Token) (s string) {
	for _, t := range tokens {
		if p.token() == t {
			s = p.text()
			p.consume()
			return s
		}
	}
	p.panic()
	return s
}

func (p *Parser) matchS(t Token, s string) string {
	if p.token() == t && p.text() == s {
		p.consume()
	} else {
		p.panic()
	}

	return s
}

func (p *Parser) panic() {
	err := fmt.Errorf(
		"parser: unexpected token '%s' in '%s' expected %s",
		p.text(),
		p.src,
		Names[p.token()],
	)
	panic(err)
}

func (p *Parser) print() {
	log.Printf("current token: %s (%s)", Names[p.token()], p.text())
}

func (p *Parser) trim(s string) string {
	return strings.Trim(s, `"`)
}
