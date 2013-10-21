package dsql

import (
	"database/sql/driver"
	"io"
)

var ComparisonOperators = map[string]string{
	"=":       "EQ",
	"<":       "LT",
	"<=":      "LE",
	">":       "GT",
	">=":      "GE",
	"like":    "BEGINS_WITH",
	"between": "BETWEEN",
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

type QueryResult struct {
	ConsumedCapacity struct {
		CapcityUnits int
		TableName    string
	}
	Count            int
	Items            []map[string]Item
	LastEvaluatedKey map[string]Value

	columns []string `json:"-"`
	row     int      `json:"-"`
}

func (q *QueryResult) Columns() []string {
	if q.Count > 0 && len(q.columns) == 0 {
		for k, _ := range q.Items[0] {
			q.columns = append(q.columns, k)
		}
	}
	return q.columns
}

func (q *QueryResult) Close() error {
	return nil
}

func (q *QueryResult) Next(dest []driver.Value) error {
	if q.row >= len(q.Items) {
		return io.EOF
	}

	item := q.Items[q.row]
	for i, c := range q.Columns() {
		dest[i] = item[c].Value()
	}

	q.row++
	return nil
}
