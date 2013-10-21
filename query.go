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
	AttributeValueList []Attribute
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

	var attrs []Attribute

	if _, ok := q.KeyConditions[exp.Identifier]; ok {
		attrs = q.KeyConditions[exp.Identifier].AttributeValueList
		attrs = append(attrs, exp.Attribute())
	} else {
		attrs = []Attribute{exp.Attribute()}
	}

	if exp.BetweenText != "" {
		attr := exp.BetweenAttribute()
		attrs = append(attrs, attr)
	}

	q.KeyConditions[exp.Identifier] = KeyCondition{
		ComparisonOperator: ComparisonOperators[exp.Operator],
		AttributeValueList: attrs,
	}
}

type Item map[string]Attribute

type QueryResult struct {
	ConsumedCapacity struct {
		CapcityUnits int
		TableName    string
	}
	Count            int
	Items            []Item
	LastEvaluatedKey Item

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
