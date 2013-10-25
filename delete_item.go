package dsql

import (
	"database/sql/driver"
	"io"
)

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

func (d DeleteItem) Rows(body io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}
