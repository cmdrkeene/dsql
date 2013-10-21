package dsql

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
