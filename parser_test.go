package dsql

import (
	"reflect"
	"testing"
)

func TestParseSelectWildcard(t *testing.T) {
	source := "select * from messages"
	expected := Query{TableName: "messages"}

	actual, err := Parse(source)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual  ", actual)
	}
}

func TestParseSelectAttributes(t *testing.T) {
	source := "select id, name from messages"
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
	}

	actual, err := Parse(source)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual  ", actual)
	}
}
