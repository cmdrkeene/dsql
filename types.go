package dsql

import "strings"

var DynamoOperators = map[string]string{
	"=":       "EQ",
	"<":       "LT",
	"<=":      "LE",
	">":       "GT",
	">=":      "GE",
	"like":    "BEGINS_WITH",
	"between": "BETWEEN",
}

var DynamoTypes = map[Token]string{
	Number: "N",
	String: "S",
}

var DefinitionTypes = map[string]string{
	"string":    "S",
	"stringset": "SS",
	"number":    "N",
	"numberset": "NS",
	"binary":    "B",
	"binaryset": "BS",
}

func NewValue(t Token, s string) Value {
	return Value{DynamoTypes[t], trim(s)}
}

type Value struct {
	Type  string
	Value string
}

type KeyCondition struct {
	ComparisonOperator string
	AttributeValueList []Value
}

type Query struct {
	TableName       string
	AttributesToGet []string
	KeyConditions   map[string]KeyCondition
}

func (q *Query) AddColumn(col string) {
	q.AttributesToGet = append(q.AttributesToGet, col)
}

func (q *Query) AddCondition(exp Expression) {
	if q.KeyConditions == nil {
		q.KeyConditions = map[string]KeyCondition{}
	}

	var values []Value
	value := Value{DynamoTypes[exp.ValueToken], exp.ValueText}

	if _, ok := q.KeyConditions[exp.Identifier]; ok {
		values = q.KeyConditions[exp.Identifier].AttributeValueList
		values = append(values, value)
	} else {
		values = []Value{value}
	}

	if exp.ValueBetweenText != "" {
		value = Value{DynamoTypes[exp.ValueToken], exp.ValueBetweenText}
		values = append(values, value)
	}

	q.KeyConditions[exp.Identifier] = KeyCondition{
		ComparisonOperator: DynamoOperators[exp.Operator],
		AttributeValueList: values,
	}
}

type PutItem struct {
	TableName string
	Item      map[string]Value
}

type UpdateItem struct {
	TableName        string
	Key              map[string]Value
	AttributeUpdates map[string]Update
}

func (u *UpdateItem) AddKey(exp Expression) {
	if u.Key == nil {
		u.Key = map[string]Value{}
	}

	u.Key[exp.Identifier] = NewValue(exp.ValueToken, exp.ValueText)
}

func (u *UpdateItem) AddUpdate(exp Expression) {
	if u.AttributeUpdates == nil {
		u.AttributeUpdates = map[string]Update{}
	}

	u.AttributeUpdates[exp.Identifier] = Update{
		Value:  NewValue(exp.ValueToken, exp.ValueText),
		Action: "PUT",
	}
}

type Update struct {
	Value  Value
	Action string // PUT, ADD, DELETE
}

type CreateTable struct {
	TableName            string
	AttributeDefinitions []AttributeDefinition
	KeySchema            []Schema
}

func (c *CreateTable) AddDefinition(d Definition) {
	c.AttributeDefinitions = append(
		c.AttributeDefinitions,
		AttributeDefinition{d.Identifier, DefinitionTypes[d.Type]},
	)

	if d.Constraint != "" {
		c.KeySchema = append(
			c.KeySchema,
			Schema{d.Identifier, strings.ToUpper(d.Constraint)},
		)
	}
}

type AttributeDefinition struct {
	AttributeName string
	AttributeType string
}

type Schema struct {
	AttributeName string
	KeyType       string
}
