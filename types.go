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
