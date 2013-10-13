package dsql

import (
	"reflect"
	"strings"
	"testing"
)

func TestLexerSimple(t *testing.T) {
	source := strings.NewReader("select * from messages")
	l := NewLexer(source)
	actual := l.Tokens()

	expected := []token{
		TokKeyword,
		TokWildcard,
		TokKeyword,
		TokId,
	}

	if reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual", actual)
	}
}

func TestLexerComplex(t *testing.T) {
	// source := "select group_id, message_id from messages where group_id = 1 and message_id > 2"
}
