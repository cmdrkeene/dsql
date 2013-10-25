package dsql

import (
	"encoding/json"
	"io"
)

type ResponseDecoder struct {
	body     io.Reader
	response Response
}

func (r *ResponseDecoder) Decode() error {
	decoder := json.NewDecoder(r.body)
	err := decoder.Decode(&r.response)
	if err != nil {
		return err
	}
	return nil
}
