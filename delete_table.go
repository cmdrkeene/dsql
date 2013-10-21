package dsql

import (
	"database/sql/driver"
	"io"
)

type DeleteTable struct {
	TableName string
}

func (d DeleteTable) Result(body io.ReadCloser) (driver.Rows, error) {
	return nil, nil
}
