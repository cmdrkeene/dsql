package dsql

import (
	"log"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	sql := "select group_id, message_id from messages where group_id = 1 and message_id > 2"
	log.Print(sql)
	src := strings.NewReader(sql)
	tokenizer := &Tokenizer{src}
	tokenizer.Tokenize()
	t.Error("test")
}
