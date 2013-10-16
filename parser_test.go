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
				AttributeValueList: []Value{Value{"N", "1"}},
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

func TestParseSelectMultipleConditions(t *testing.T) {
	source := `select id, name from messages where id = 1 AND name = "a"`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Value{Value{"N", "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Value{Value{"S", "a"}},
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

func TestParseSelectMultipleConditionsOr(t *testing.T) {
	source := `select id, name from messages where id = 1 AND name = "a" OR name = "b"`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Value{Value{"N", "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Value{
					Value{"S", "a"},
					Value{"S", "b"},
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

func TestParseSelectMultipleConditionsBetween(t *testing.T) {
	source := `select id, name from messages where id = 1 AND name BETWEEN("a", "z")`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Value{Value{"N", "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "BETWEEN",
				AttributeValueList: []Value{
					Value{"S", "a"},
					Value{"S", "z"},
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

func TestParseInsert(t *testing.T) {
	source := `insert into messages (id, name) values (1, "a")`
	expected := PutItem{
		TableName: "messages",
		Item: map[string]Value{
			"id":   Value{"N", "1"},
			"name": Value{"S", "a"},
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
