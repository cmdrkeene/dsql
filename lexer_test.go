package dsql

import (
	"strings"

	"testing"
)

func TestTokenize(t *testing.T) {
	source := "select group_id, message_id from messages_v4_production where group_id = 1 and message_id > 2"
	l := NewLexer(strings.NewReader(source))
	t.Error(l)
}
