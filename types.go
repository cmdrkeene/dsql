// Declare types and conversion maps
package dsql

import (
	"database/sql/driver"
	"encoding/json"
	"io"

	"strconv"
	"strings"
)

type Operation string
type Request interface {
	Rows(io.ReadCloser) (driver.Rows, error)
}

var ComparisonOperators = map[string]string{
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
		ComparisonOperator: ComparisonOperators[exp.Operator],
		AttributeValueList: values,
	}
}

func (q Query) Rows(r io.ReadCloser) (driver.Rows, error) {
	defer r.Close()
	dec := json.NewDecoder(r)
	result := QueryResult{}
	err := dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (qr QueryResult) Columns() []string {
	return []string{"ass", "ass", "ass"}
}

func (qr QueryResult) Close() error {
	return nil
}

func (qr QueryResult) Next(dest []driver.Value) error {
	return io.EOF
}

type PutItem struct {
	TableName string
	Item      map[string]Value
}

func (p PutItem) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}

type UpdateItem struct {
	TableName        string
	Key              map[string]Value
	AttributeUpdates map[string]Update
}

func (u UpdateItem) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}

func (u UpdateItem) AddKey(exp Expression) {
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

func (u Update) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
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

func (c CreateTable) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
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

func (d DeleteItem) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
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

func (d DeleteTable) Rows(r io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}

func operation(r Request) Operation {
	switch r.(type) {
	case Query:
		return Operation("DynamoDB_20120810.Query")
	}
	return ""
}
