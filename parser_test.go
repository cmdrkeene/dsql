package dsql

import (
	"reflect"
	"testing"
)

func TestParseError(t *testing.T) {
	var err error

	_, err = Parse("oh hai frenz")
	if err != ErrSyntax {
		t.Error("unexpected ", err)
	}

	_, err = Parse("selects")
	if err != ErrSyntax {
		t.Error("unexpected ", err)
	}

	_, err = Parse("create table foo(")
	if err == nil {
		t.Error("unexpected ", err)
	}
}

func TestParseSelectWildcard(t *testing.T) {
	source := "select * from messages;"
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

func TestParseSelectWithoutSemicolon(t *testing.T) {
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
	source := "select id, name from messages;"
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
	source := "select id, name from messages where id = 1;"
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{N: "1"}},
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
	source := `select id, name from messages where id = 1 AND name = "a";`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{N: "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{S: "a"}},
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
	source := `select id, name from messages where id = 1 AND name = "a" OR name = "b";`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{N: "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{S: "a"}, Attribute{S: "b"}},
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
	source := `select id, name from messages where id = 1 AND name BETWEEN("a", "z");`
	expected := Query{
		TableName:       "messages",
		AttributesToGet: []string{"id", "name"},
		KeyConditions: map[string]KeyCondition{
			"id": KeyCondition{
				ComparisonOperator: "EQ",
				AttributeValueList: []Attribute{Attribute{N: "1"}},
			},
			"name": KeyCondition{
				ComparisonOperator: "BETWEEN",
				AttributeValueList: []Attribute{Attribute{S: "a"}, Attribute{S: "z"}},
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
	source := `insert into messages (id, name) values (1, "a");`
	expected := PutItem{
		TableName: "messages",
		Item: Item{
			"id":   Attribute{N: "1"},
			"name": Attribute{S: "a"},
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

func TestParseCreate(t *testing.T) {
	source := `create table messages (group string hash, id number range);`
	expected := CreateTable{
		TableName: "messages",
		AttributeDefinitions: []AttributeDefinition{
			AttributeDefinition{
				AttributeName: "group",
				AttributeType: "S",
			},
			AttributeDefinition{
				AttributeName: "id",
				AttributeType: "N",
			},
		},
		KeySchema: []Schema{
			Schema{
				AttributeName: "group",
				KeyType:       "HASH",
			},
			Schema{
				AttributeName: "id",
				KeyType:       "RANGE",
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

func TestParseCreateThroughput(t *testing.T) {
	source := `create table messages (group string hash, id number range) with (read = 10, write = 5);`
	expected := CreateTable{
		TableName: "messages",
		AttributeDefinitions: []AttributeDefinition{
			AttributeDefinition{
				AttributeName: "group",
				AttributeType: "S",
			},
			AttributeDefinition{
				AttributeName: "id",
				AttributeType: "N",
			},
		},
		KeySchema: []Schema{
			Schema{
				AttributeName: "group",
				KeyType:       "HASH",
			},
			Schema{
				AttributeName: "id",
				KeyType:       "RANGE",
			},
		},
	}
	expected.ProvisionedThroughput.ReadCapacityUnits = 10
	expected.ProvisionedThroughput.WriteCapacityUnits = 5

	actual, err := Parse(source)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual  ", actual)
	}
}

func TestDropTable(t *testing.T) {
	source := `drop table messages;`
	expected := DeleteTable{TableName: "messages"}

	actual, err := Parse(source)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected", expected)
		t.Error("actual  ", actual)
	}
}

func TestDeleteItem(t *testing.T) {
	source := `delete from messages where id = 1 AND name = "a";`
	expected := DeleteItem{
		TableName: "messages",
		Key: Item{
			"id":   Attribute{N: "1"},
			"name": Attribute{S: "a"},
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
