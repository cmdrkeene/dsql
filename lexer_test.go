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
		TokEOF,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual", actual)
	}
}

func TestLexerComplex(t *testing.T) {
	source := strings.NewReader(`
    SELECT group_id, message_id
    FROM messages
    WHERE group_id = "1" AND message_id > 2
  `)
	l := NewLexer(source)
	actual := l.Tokens()

	expected := []token{
		TokKeyword, // select
		TokId,      // group_id
		TokComma,   // ,
		TokId,      // message_id
		TokKeyword, // from
		TokId,      // messages
		TokKeyword, // where
		TokId,      // group_id
		TokOp,      // =
		TokString,  // "1"
		TokKeyword, // and
		TokId,      // message_id
		TokOp,      // >
		TokNumber,  // 2
		TokEOF,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", l.Names(expected))
		t.Error("actual  ", l.Names(actual))
	}
}
