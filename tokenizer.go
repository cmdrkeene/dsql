package dsql

import (
	"io"
	"log"
	"text/scanner"
)

type Tokenizer struct {
	src io.Reader
}

func (t *Tokenizer) Tokenize() {
	var s scanner.Scanner

	s.Init(t.src)
	tok := s.Scan()
	log.Print(scanner.TokenString(tok), ":", s.TokenText())
	for tok != scanner.EOF {
		tok = s.Scan()
		log.Print(scanner.TokenString(tok), ":", s.TokenText())
	}
	log.Print(s.TokenText())
}
