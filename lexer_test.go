package dsql

import (
	"reflect"
	"strings"
	"testing"
)

type Case struct {
	source string
	tokens []token
}

func TestLexerCases(t *testing.T) {
	cases := []Case{
		Case{
			`select * from table`,
			[]token{K, W, K, I, E},
		},
		Case{
			`SELECT id, name FROM table WHERE id = "1" AND name > 2`,
			[]token{K, I, C, I, K, I, K, I, O, S, K, I, O, N, E},
		},
	}

	var l *Lexer
	var expected, actual []token

	for _, c := range cases {
		l = NewLexer(strings.NewReader(c.source))
		actual = l.Tokens()
		expected = c.tokens

		if !reflect.DeepEqual(actual, expected) {
			t.Error("source", c.source)
			t.Error("expected", l.Names(expected))
			t.Error("actual  ", l.Names(actual))
		}
	}
}
