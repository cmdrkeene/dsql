package dsql

type Token rune

const (
	Keywords    = "^(select|insert|create|update|delete|drop|from|where|set|limit|order|by|asc|desc|into|values|table|with)$"
	Types       = "^(number|numberset|string|stringset)$"
	Constraints = "^(hash|range|index|all|projection)$"
	Operators   = "^(=|>|>=|<|<=|like|and|or|between)$"
)

const (
	Keyword Token = iota
	Identifier
	Constraint
	Wildcard
	Type
	String
	Number
	Operator
	Comma
	LeftParen
	RightParen
	EOF
	Unknown
)

var Names = map[Token]string{
	Keyword:    "Keyword",
	Identifier: "Identifier",
	Constraint: "Constraint",
	Wildcard:   "Wildcard",
	Type:       "Type",
	String:     "String",
	Number:     "Number",
	Operator:   "Operator",
	Comma:      ",",
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
	X = Constraint
	W = Wildcard
	T = Type
	N = Number
	S = String
	O = Operator
	C = Comma
	L = LeftParen
	R = RightParen
	E = EOF
	U = Unknown
)

var Symbols = map[Token]string{
	Keyword:    "K",
	Identifier: "I",
	Constraint: "X",
	Wildcard:   "W",
	Type:       "T",
	String:     "S",
	Number:     "N",
	Operator:   "O",
	Comma:      "C",
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
