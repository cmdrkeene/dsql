package dsql

type PutItem struct {
	TableName string
	Item      map[string]Value
}
