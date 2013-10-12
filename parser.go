package dsql

import "io"

type Parser struct {
	lexer *Lexer
	tok   token
}

func NewParser(source io.Reader) *Parser {
	return &Parser{lexer: NewLexer(source)}
}

func (p *Parser) Parse() (err error) {
	for p.tok != TokEOF {
		p.tok, err = p.lexer.Next()

		switch p.tok {
		case TokSelect:
			return p._select()
		default:
			return ErrUnexpectedToken
		}
	}

	return nil
}

func (p *Parser) _select() error {
	return nil
}
