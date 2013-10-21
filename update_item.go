package dsql

type UpdateItem struct {
	TableName        string
	Key              Item
	AttributeUpdates map[string]Update
}

func (u *UpdateItem) AddKey(exp Expression) {
	if u.Key == nil {
		u.Key = Item{}
	}

	u.Key[exp.Identifier] = exp.Attribute()
}

func (u *UpdateItem) AddUpdate(exp Expression) {
	if u.AttributeUpdates == nil {
		u.AttributeUpdates = map[string]Update{}
	}

	u.AttributeUpdates[exp.Identifier] = Update{
		Value:  exp.Attribute(),
		Action: "PUT",
	}
}

type Update struct {
	Value  Attribute
	Action string // PUT, ADD, DELETE
}
