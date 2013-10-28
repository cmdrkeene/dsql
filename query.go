package dsql

import (
	"database/sql/driver"
	"io"
)

type Item map[string]Attribute

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
	TableName        string
	AttributesToGet  []string
	KeyConditions    map[string]KeyCondition
	ScanIndexForward bool
	Limit            int
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

func (q *Query) Rows(body io.ReadCloser) (driver.Rows, error) {
	res := &QueryResponse{}
	dec := ResponseDecoder{body, res}
	err := dec.Decode()
	if err != nil {
		return nil, err
	}

	if q.AttributesToGet != nil {
		res.cols = q.AttributesToGet
	}

	return &Rows{cols: res.Columns(), values: res.Values()}, nil
}
