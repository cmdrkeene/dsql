package dsql

import (
	"encoding/json"
	"github.com/groupme/dynamo"
	"strconv"
	"strings"
)

type Request interface{}
type Result interface{}

var QueryOperators = map[string]string{
	"=":       "EQ",
	"<":       "LT",
	"<=":      "LE",
	">":       "GT",
	">=":      "GE",
	"like":    "BEGINS_WITH",
	"between": "BETWEEN",
}

var TokenTypes = map[Token]string{
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
	return Value{TokenTypes[t], trim(s)}
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
	result          QueryResult
}

func (q *Query) AddColumn(col string) {
	q.AttributesToGet = append(q.AttributesToGet, col)
}

func (q *Query) AddCondition(exp Expression) {
	if q.KeyConditions == nil {
		q.KeyConditions = map[string]KeyCondition{}
	}

	var values []Value
	value := Value{TokenTypes[exp.ValueToken], exp.ValueText}

	if _, ok := q.KeyConditions[exp.Identifier]; ok {
		values = q.KeyConditions[exp.Identifier].AttributeValueList
		values = append(values, value)
	} else {
		values = []Value{value}
	}

	if exp.ValueBetweenText != "" {
		value = Value{TokenTypes[exp.ValueToken], exp.ValueBetweenText}
		values = append(values, value)
	}

	q.KeyConditions[exp.Identifier] = KeyCondition{
		ComparisonOperator: QueryOperators[exp.Operator],
		AttributeValueList: values,
	}
}

type QueryResult struct {
	ConsumedCapacity struct {
		CapcityUnits int
		TableName    string
	}
	Count            int
	Items            map[string]Value
	LastEvaluatedKey map[string]Value
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
	TableName             string
	AttributeDefinitions  []AttributeDefinition
	KeySchema             []Schema
	ProvisionedThroughput struct {
		ReadCapacityUnits  int
		WriteCapacityUnits int
	}
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

func (c *CreateTable) AddThroughput(exp Expression) {
	units, err := strconv.Atoi(exp.ValueText)
	if err != nil {
		panic("throughput must be an integer")
	}

	switch exp.Identifier {
	case "read":
		c.ProvisionedThroughput.ReadCapacityUnits = units
	case "write":
		c.ProvisionedThroughput.WriteCapacityUnits = units
	default:
		panic("unknown create table parameter (expected read or write)")
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

type DeleteItem struct {
	TableName string
	Key       map[string]Value
}

func (d *DeleteItem) AddKey(exp Expression) {
	if d.Key == nil {
		d.Key = map[string]Value{}
	}
	d.Key[exp.Identifier] = Value{TokenTypes[exp.ValueToken], exp.ValueText}
}

type DeleteTable struct {
	TableName string
}

func marshal(r Request) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func operation(r Request) string {
	switch r.(type) {
	case Query:
		return dynamo.OpQuery
	}
	return ""
}

func result(r Request) Result {
	switch r.(type) {
	case Query:
		return QueryResult{}
	}
	return nil
}
