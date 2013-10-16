package dsql

type KeyCondition struct {
	ComparisonOperator string
	AttributeValueList []map[string]string
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
	q.KeyConditions = map[string]KeyCondition{
		exp.Identifier: KeyCondition{
			ComparisonOperator: DynamoOperators[exp.Operator],
			AttributeValueList: []map[string]string{
				map[string]string{
					DynamoTypes[exp.ValueToken]: exp.ValueText,
				},
			},
		},
	}
}

var DynamoOperators = map[string]string{
	"=":    "EQ",
	"<":    "LT",
	"<=":   "LE",
	">":    "GT",
	">=":   "GE",
	"like": "BEGINS_WITH",
}

var DynamoTypes = map[Token]string{
	Number: "N",
	String: "S",
}
