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

// recover from unexpected eof
// recover from unexpected token
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

	if !p.expect(Wildcard, Identifier) {
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

	if !p.matchKeyword("from") {
		return nil
	}

	p.next()
	if !p.match(Identifier) {
		return nil
	}

	req.TableName = p.text()

	if p.next() == EOF {
		return req
	}

	// conditions
	// WHERE ((Ident Operator (Number|String)))(Comma)?)+ (EOF|Keyword)

	if !p.matchKeyword("where") {
		return nil
	}

	p.next()
	for p.match(Identifier, Comma) {
		if p.peek() == Identifier {
			attr := p.text()
			kc := KeyCondition{}

			// operator
			p.next()
			if !p.expect(Operator) {
				return nil
			}

			kc.ComparisonOperator = DynamoOperators[p.text()]

			// value
			p.next()
			if !p.expect(Number, String) {
				return nil
			}
			v := map[string]string{DynamoTypes[p.peek()]: p.text()}
			kc.AttributeValueList = append(kc.AttributeValueList, v)
			req.KeyConditions[attr] = kc
		}
		p.next()
	}

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
	return false
}

func (p *Parser) expect(tokens ...Token) bool {
	if p.match(tokens...) {
		return true
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

func (p *Parser) matchKeyword(s string) bool {
	if p.peek() == Keyword && p.text() == s {
		return true
	}
	p.errUnexpected()
	return false
}

func (p *Parser) errUnexpected() {
	p.err = fmt.Errorf("parser: unexpected token '%s' in '%s' expected %s", p.text(), p.src, Names[p.peek()])
}
