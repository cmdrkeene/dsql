package dsql

import (
	"database/sql/driver"
	"io"
)

type PutItem struct {
	TableName string
	Item      Item
}

func (p PutItem) Rows(body io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}
