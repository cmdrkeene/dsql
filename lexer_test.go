package dsql

import (
	"reflect"
	"strings"
	"testing"
)

type Case struct {
	source string
	tokens []Token
}

func TestLexerCases(t *testing.T) {
	cases := []Case{
		Case{
			`SELECT * FROM users`,
			[]Token{K, W, K, I, E},
		},
		Case{
			`SELECT * FROM users LIMIT 10`,
			[]Token{K, W, K, I, K, N, E},
		},
		Case{
			`SELECT * FROM users LIMIT 10 ORDER BY id ASC`,
			[]Token{K, W, K, I, K, N, K, K, I, K, E},
		},
		Case{
			`SELECT id, name FROM users`,
			[]Token{K, I, C, I, K, I, E},
		},
		Case{
			`SELECT id, name FROM users WHERE id = "1" AND name > 2`,
			[]Token{K, I, C, I, K, I, K, I, O, S, K, I, O, N, E},
		},
		Case{
			`SELECT id, name FROM users WHERE id = 1 AND name > "A"`,
			[]Token{K, I, C, I, K, I, K, I, O, N, K, I, O, S, E},
		},
		Case{
			`INSERT INTO users (id, name) VALUES (1, "A")`,
			[]Token{K, K, I, L, I, C, I, R, K, L, N, C, S, R, E},
		},
		Case{
			`UPDATE users SET name = "B" WHERE name = "A"`,
			[]Token{K, I, K, I, O, S, K, I, O, S, E},
		},
		Case{
			`DELETE FROM users WHERE name = "A"`,
			[]Token{K, K, I, K, I, O, S, E},
		},
		Case{
			`
			CREATE TABLE messages (
			    GroupId STRING HASH KEY,
			    MessageId NUMBER RANGE KEY,
			    UserId NUMBER,
			    Text STRING
			);`,
			[]Token{K, K, I, L, I, T, T, T, C, I, T, T, T, C, I, T, C, I, T, R, M, E},
		},
	}

	var l *Lexer
	var expected, actual []Token

	for _, c := range cases {
		l = NewLexer(strings.NewReader(c.source))
		actual = l.tokens
		expected = c.tokens

		if !reflect.DeepEqual(actual, expected) {
			t.Error("source", c.source)
			t.Error("expected", symbols(expected))
			t.Error("actual  ", symbols(actual))
		}
	}
}

func TestPeekNextText(t *testing.T) {
	source := `SELECT name FROM users`

	l := NewLexer(strings.NewReader(source))

	if l.Peek() != Keyword {
		t.Error("expected keyword", Names[l.Peek()])
	}

	if l.Text() != "select" {
		t.Error("expected select", l.Text())
	}

	if l.Next() != Identifier {
		t.Error("expected identifier", Names[l.Peek()])
	}

	if l.Next() != Keyword {
		t.Error("expected keyword", Names[l.Peek()])
	}

	if l.Next() != Identifier {
		t.Error("expected identifier", Names[l.Peek()])
	}

	if l.Next() != EOF {
		t.Error("expected EOF", Names[l.Peek()])
	}
}
