// Might want to annotate with special info: consumed units, etc.
package dsql

import (
	"database/sql/driver"
	"io"
)

// driver implentation of driver.Rows interface
type Rows struct {
	values [][]driver.Value
	cols   []string
	idx    int
}

func (r *Rows) Columns() []string {
	return r.cols
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	if r.idx == len(r.values) {
		return io.EOF
	}
	for i, v := range r.values[r.idx] {
		dest[i] = v
	}
	r.idx++
	return nil
}
