package dsql

type Token rune

const (
	Keywords  = "select|insert|create|update|delete|from|where|and|set|limit|order|by|asc|desc|into|values|table"
	Types     = "number|string|hash|range|key"
	Operators = "=|>|>=|<|<=|like"
)

const (
	Keyword Token = iota
	Identifier
	Wildcard
	Type
	String
	Number
	Operator
	Comma
	Semicolon
	LeftParen
	RightParen
	EOF
	Unknown
)

var Names = map[Token]string{
	Keyword:    "Keyword",
	Identifier: "Identifier",
	Wildcard:   "Wildcard",
	Type:       "Type",
	String:     "String",
	Number:     "Number",
	Operator:   "Operator",
	Comma:      ",",
	Semicolon:  ";",
	LeftParen:  "(",
	RightParen: ")",
	EOF:        "EOF",
	Unknown:    "Unknown",
}

func names(tokens []Token) (n []string) {
	for _, t := range tokens {
		n = append(n, Names[t])
	}
	return n
}

const (
	K = Keyword
	I = Identifier
	W = Wildcard
	T = Type
	N = Number
	S = String
	O = Operator
	C = Comma
	M = Semicolon
	L = LeftParen
	R = RightParen
	E = EOF
	U = Unknown
)

var Symbols = map[Token]string{
	Keyword:    "K",
	Identifier: "I",
	Wildcard:   "W",
	Type:       "T",
	String:     "S",
	Number:     "N",
	Operator:   "O",
	Comma:      "C",
	Semicolon:  "M",
	LeftParen:  "L",
	RightParen: "R",
	EOF:        "E",
	Unknown:    "U",
}

func symbols(tokens []Token) (s []string) {
	for _, t := range tokens {
		s = append(s, Symbols[t])
	}
	return s
}
