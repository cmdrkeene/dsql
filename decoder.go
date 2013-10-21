package dsql

import (
	"database/sql/driver"
	"encoding/json"
	"io"
)

func decode(req Request, body io.ReadCloser) (driver.Rows, error) {
	res := QueryResult{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
