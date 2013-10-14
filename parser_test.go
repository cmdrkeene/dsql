package dsql

import (
	"reflect"
	"testing"
)

func TestParseSelect(t *testing.T) {
	source := "select * from messages"
	expected := Query{
		TableName: "messages",
		Select:    "ALL_ATTRIBUTES",
	}

	actual, err := Parse(source)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("source", source)
		t.Error("expected", expected)
		t.Error("actual  ", actual)
	}
}
