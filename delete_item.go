package dsql

type DeleteItem struct {
	TableName string
	Key       Item
}

func (d *DeleteItem) AddKey(exp Expression) {
	if d.Key == nil {
		d.Key = Item{}
	}
	d.Key[exp.Identifier] = exp.Attribute()
}
