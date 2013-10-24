// Convert tokens into requests
package dsql

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var ErrSyntax = errors.New("parser: syntax error")

func Parse(source string) (Request, error) {
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

func (p *Parser) Parse() (req Request, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	switch p.token() {
	case Keyword:
		switch p.text() {
		case "select":
			req = p.Select()
		case "insert":
			req = p.Insert()
		case "update":
			req = p.Update()
		case "create":
			req = p.Create()
		case "delete":
			req = p.Delete()
		case "drop":
			req = p.Drop()
		default:
			p.err = ErrSyntax
		}
	default:
		p.err = ErrSyntax
	}

	return req, p.err
}

func (p *Parser) Select() Request {
	query := Query{}
	p.columns(&query)
	p.from(&query)
	p.limit(&query)
	p.where(&query)
	p.limit(&query)
	return query
}

func (p *Parser) from(query *Query) {
	p.matchS(Keyword, "from")
	query.TableName = p.match(Identifier)
}

func (p *Parser) columns(query *Query) {
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
}

func (p *Parser) limit(query *Query) {
	if p.token() == Keyword && p.text() == "limit" {
		p.consume()
		limit, err := strconv.Atoi(p.match(Number))
		if err != nil {
			panic(err)
		}
		query.Limit = limit
	}
}

func (p *Parser) where(query *Query) {
	if p.token() == Keyword && p.text() == "where" {
		p.consume()
		query.AddCondition(p.expr())
		for p.token() == Operator {
			p.match(Operator)
			query.AddCondition(p.expr())
		}
	}
}

func (p *Parser) Insert() Request {
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

	attrs := []Attribute{p.attribute()}
	p.match(String, Number)
	for p.token() == Comma {
		p.match(Comma)
		attr := p.attribute()
		p.match(String, Number)
		attrs = append(attrs, attr)
	}

	p.match(RightParen)

	item := Item{}

	for i, col := range columns {
		item[col] = attrs[i]
	}

	return PutItem{TableName: table, Item: item}
}

func (p *Parser) Update() Request {
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

func (p *Parser) Create() Request {
	create := CreateTable{}

	p.matchS(Keyword, "create")
	p.matchS(Keyword, "table")
	create.TableName = p.match(Identifier)
	p.match(LeftParen)

	create.AddDefinition(p.definition())
	for p.token() == Comma {
		p.match(Comma)
		create.AddDefinition(p.definition())
	}

	p.match(RightParen)

	if p.token() == Keyword {
		p.matchS(Keyword, "with")
		p.match(LeftParen)
		create.AddThroughput(p.expr())
		for p.token() == Comma {
			p.match(Comma)
			create.AddThroughput(p.expr())
		}
		p.match(RightParen)
	}

	return create
}

func (p *Parser) Delete() Request {
	p.matchS(Keyword, "delete")
	p.matchS(Keyword, "from")
	deleteItem := DeleteItem{TableName: p.match(Identifier)}
	p.matchS(Keyword, "where")
	deleteItem.AddKey(p.expr())
	for p.token() == Operator {
		p.matchS(Operator, "and")
		deleteItem.AddKey(p.expr())
	}

	return deleteItem
}

func (p *Parser) Drop() Request {
	p.matchS(Keyword, "drop")
	p.matchS(Keyword, "table")
	return DeleteTable{p.match(Identifier)}
}

func (p *Parser) consume() Token {
	return p.lex.Next()
}

func (p *Parser) token() Token {
	return p.lex.Peek()
}

func (p *Parser) text() string {
	return strings.Trim(p.lex.Text(), `"`)
}

func (p *Parser) attribute() (a Attribute) {
	switch p.token() {
	case String:
		a.S = p.text()
	case Number:
		a.N = p.text()
	default:
		p.print()
		panic("unknown attribute")
	}
	return a
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
		"parser: unexpected token '%s' (%s) in '%s'",
		p.text(),
		Names[p.token()],
		p.src,
	)
	panic(err)
}

func (p *Parser) print() {
	log.Printf("current token: %s (%s)", Names[p.token()], p.text())
}

type Expression struct {
	Identifier  string
	Operator    string
	Token       Token
	Text        string
	BetweenText string
}

func (exp *Expression) Attribute() (a Attribute) {
	return exp.newAttr(exp.Text)
}

func (exp *Expression) BetweenAttribute() (a Attribute) {
	return exp.newAttr(exp.BetweenText)
}

func (exp *Expression) newAttr(s string) (a Attribute) {
	switch exp.Token {
	case String:
		a = Attribute{S: s}
	case Number:
		a = Attribute{N: s}
	}
	return a
}

func (p *Parser) expr() (exp Expression) {
	exp.Identifier = p.match(Identifier)

	if p.text() == "between" {
		exp.Operator = p.match(Operator)
		p.match(LeftParen)
		exp.Token = p.token()
		exp.Text = p.match(String, Number)
		p.match(Comma)
		exp.BetweenText = p.match(String, Number)
		p.match(RightParen)
	} else {
		exp.Operator = p.match(Operator)
		exp.Token = p.token()
		exp.Text = p.match(String, Number)
	}
	return exp
}

type Definition struct {
	Identifier string
	Type       string
	Constraint string
}

func (p *Parser) definition() (def Definition) {
	def.Identifier = p.match(Identifier)
	def.Type = p.match(Type)
	if p.token() == Constraint {
		def.Constraint = p.match(Constraint)
	}
	return def
}
