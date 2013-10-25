package dsql

import (
	"database/sql/driver"
	"io"
	"log"
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

func (q *Query) Rows(body io.ReadCloser) (rows driver.Rows, err error) {
	dec := ResponseDecoder{body, &QueryResponse{}}
	err = dec.Decode()
	if err != nil {
		return nil, err
	}

	rows = NewRows(dec.response)

	log.Print("response: ", dec.response)
	log.Print("rows: ", rows)

	return rows, nil
}
