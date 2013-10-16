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

func TestParseSelectColumns(t *testing.T) {
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

func TestParseSelectSingleCondition(t *testing.T) {
	source := "select id, name from messages where id = 1"
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []map[string]string{
					map[string]string{"N": "1"},
				},
			},
		},
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

// func TestParseSelectMultipleConditions(t *testing.T) {
// 	source := `select id, name from messages where id = 1 AND name = "a"`
// 	expected := Query{
// 		TableName:       "messages",
// 		AttributesToGet: []string{"id", "name"},
// 		KeyConditions: map[string]KeyCondition{
// 			"id": KeyCondition{
// 				ComparisonOperator: "EQ",
// 				AttributeValueList: []map[string]string{
// 					map[string]string{"N": "1"},
// 				},
// 			},
// 			"name": KeyCondition{
// 				ComparisonOperator: "EQ",
// 				AttributeValueList: []map[string]string{
// 					map[string]string{"S": "a"},
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
