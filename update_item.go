package dsql

type UpdateItem struct {
	TableName        string
	Key              map[string]Value
	AttributeUpdates map[string]Update
}

func (u *UpdateItem) AddKey(exp Expression) {
	if u.Key == nil {
		u.Key = map[string]Value{}
	}

	u.Key[exp.Identifier] = NewValue(exp.ValueToken, exp.ValueText)
}

func (u *UpdateItem) AddUpdate(exp Expression) {
	if u.AttributeUpdates == nil {
		u.AttributeUpdates = map[string]Update{}
	}

	u.AttributeUpdates[exp.Identifier] = Update{
		Value:  NewValue(exp.ValueToken, exp.ValueText),
		Action: "PUT",
	}
}

type Update struct {
	Value  Value
	Action string // PUT, ADD, DELETE
}
