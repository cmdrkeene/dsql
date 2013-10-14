package dsql

// import "io"

// var ErrUnexpectedToken = errors.New("lexer: unexpected token")

// type QueryRequest struct {
// 	TableName       string
// 	AttributesToGet []string
// }

// type Parser struct {
// 	lexer *Lexer
// 	tok   token
// }

// func NewParser(source io.Reader) *Parser {
// 	return &Parser{lexer: NewLexer(source)}
// }

// func (p *Parser) Parse() (request interface{}, err error) {
// 	// recover from panic
// 	p.next()
// 	for p.tok != TokEOF {
// 		switch p.tok {
// 		case TokSelect:
// 			return p._select()
// 		default:
// 			return ErrUnexpectedToken
// 		}
// 	}

// 	return nil
// }

// // list
// // assign

// func (p *Parser) _select() (req interface{}, err error) {
// 	// req = &QueryRequest{}

// 	// p.match(TokId)
// 	// p.match(TokFrom)
// 	// p.match(TokId)
// 	// p.match(TokWhere)
// 	// p.match(TokId)
// 	// p.match(TokOp)
// 	// p.match(TokId)

// 	return req, nil
// }

// func (p *Parser) next() {
// 	p.tok, err = p.lexer.Next()
// }

// func (p *Parser) match(t token) {
// 	p.next()
// 	if p.tok != t {
// 		panic("unexpected token", p.tok)
// 	}
// }
