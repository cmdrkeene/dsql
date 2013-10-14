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

// func TestParseSelectAttributesConditions(t *testing.T) {
// 	source := "select id, name from messages where id = 1"
// 	expected := Query{
// 		TableName:       "messages",
// 		AttributesToGet: []string{"id", "name"},
// 		KeyConditions: map[string]KeyCondition{
// 			"messages": KeyCondition{
// 				ComparisonOperator: "EQ",
// 				AttributeValueList: []map[string]string{
// 					map[string]string{"N": "1"},
// 				},
// 			},
// 		},
// 	}

// 	actual, err := Parse(source)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if !reflect.DeepEqual(actual, expected) {
// 		t.Error("expected", expected)
// 		t.Error("actual  ", actual)
// 	}
// }
