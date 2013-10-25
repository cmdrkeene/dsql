package dsql

import (
	"io"

	"database/sql/driver"
)

func NewRows(res Response) *Rows {
	return &Rows{
		response: res,
		cols:     res.Columns(),
		values:   res.Values(),
	}
}

// driver implentation of driver.Rows interface
// also provides access to the raw response object for advanced applications
type Rows struct {
	response interface{}

	// value representations of underlying response
	values [][]driver.Value

	// names of columns in response
	// this is problematic due to the schemaless nature of dynamo
	// driver wants fixed-width scan operations but items are variable-width
	// consider preprocessing and properly accepting nil gaps
	cols []string

	// current position in Next() operations
	idx int
}

// There might be a conventional name for this
func (r *Rows) Raw() interface{} {
	return r.response
}

func (r *Rows) Columns() []string {
	return r.cols
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	if r.idx > len(r.values) {
		return io.EOF
	}

	for i, v := range r.values[r.idx] {
		if v != nil {
			dest[i] = v
		}
	}

	return nil
}
