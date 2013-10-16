package dsql

type AttributeValue struct {
	Type  string
	Value string
}

type KeyCondition struct {
	ComparisonOperator string
	AttributeValueList []AttributeValue
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

	var list []AttributeValue
	value := AttributeValue{DynamoTypes[exp.ValueToken], exp.ValueText}

	if _, ok := q.KeyConditions[exp.Identifier]; ok {
		list = q.KeyConditions[exp.Identifier].AttributeValueList
		list = append(list, value)
	} else {
		list = []AttributeValue{value}
	}

	if exp.ValueBetweenText != "" {
		value = AttributeValue{DynamoTypes[exp.ValueToken], exp.ValueBetweenText}
		list = append(list, value)
	}

	q.KeyConditions[exp.Identifier] = KeyCondition{
		ComparisonOperator: DynamoOperators[exp.Operator],
		AttributeValueList: list,
	}
}

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
