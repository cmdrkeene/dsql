package dsql

import (
	"database/sql/driver"
	"io"
)

type Update struct {
	Value  Attribute
	Action string // PUT, ADD, DELETE
}

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

func (u UpdateItem) Result(body io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}
