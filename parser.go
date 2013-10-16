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

func (p *Parser) Parse() (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
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
		case "update":
			req = p.Update()
		}
	}
	return req, p.err
}

func (p *Parser) Select() interface{} {
	query := Query{}

	p.match(Keyword)

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

	p.matchS(Keyword, "where")

	query.AddCondition(p.expr())

	for p.token() == Operator {
		p.match(Operator)
		query.AddCondition(p.expr())
	}

	return query
}

type Expression struct {
	Identifier       string
	Operator         string
	ValueToken       Token
	ValueText        string
	ValueBetweenText string
}

func (p *Parser) expr() (exp Expression) {
	exp.Identifier = p.match(Identifier)

	if p.text() == "between" {
		exp.Operator = p.match(Operator)
		p.match(LeftParen)
		exp.ValueToken = p.token()
		exp.ValueText = trim(p.match(String, Number))
		p.match(Comma)
		exp.ValueBetweenText = trim(p.match(String, Number))
		p.match(RightParen)
	} else {
		exp.Operator = p.match(Operator)
		exp.ValueToken = p.token()
		exp.ValueText = trim(p.match(String, Number))
	}
	return exp
}

// insert into name (id, name) values (1, "a")
// TODO insert into name (id, name) values (1, "a", 2, "b")
func (p *Parser) Insert() interface{} {
	p.matchS(Keyword, "insert")
	p.matchS(Keyword, "into")

	table := p.match(Identifier)

	p.match(LeftParen)
	columns := []string{p.match(Identifier)}
	for p.token() == Comma {
		p.match(Comma)
		col := p.text()
		p.match(Identifier)
		columns = append(columns, col)
	}
	p.match(RightParen)

	p.matchS(Keyword, "values")
	p.match(LeftParen)
	values := []Value{NewValue(p.token(), p.text())}
	p.match(String, Number)
	for p.token() == Comma {
		p.match(Comma)
		value := NewValue(p.token(), p.text())
		p.match(String, Number)
		values = append(values, value)
	}
	p.match(RightParen)

	item := map[string]Value{}

	for i, col := range columns {
		item[col] = values[i]
	}

	putItem := PutItem{
		TableName: table,
		Item:      item,
	}
	return putItem
}

// update messages set name = "b" where id = 1
func (p *Parser) Update() interface{} {
	update := UpdateItem{}

	p.matchS(Keyword, "update")
	update.TableName = p.match(Identifier)

	p.matchS(Keyword, "set")

	update.AddUpdate(p.expr())
	for p.token() == Comma {
		p.match(Comma)
		update.AddUpdate(p.expr())
	}

	p.matchS(Keyword, "where")

	update.AddKey(p.expr())
	for p.token() == Operator {
		p.match(Operator)
		update.AddKey(p.expr())
	}

	return update
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
