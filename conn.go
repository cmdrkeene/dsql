// database/sql implementation
package dsql

import (
	"database/sql"
	"database/sql/driver"
)

func init() {
	sql.Register("dynamodb", &drv{})
}

type drv struct{}

func (d *drv) Open(name string) (driver.Conn, error) {
	return Open(name)
}

func Open(name string) (driver.Conn, error) {
	return &conn{cl: GetClient(name)}, nil
}

type conn struct {
	cl Client
}

func (cn *conn) Prepare(q string) (driver.Stmt, error) {
	return nil, nil
}

func (cn *conn) Close() error {
	return nil
}

func (cn *conn) Begin() (driver.Tx, error) {
	return cn, nil
}

func (cn *conn) Commit() error {
	return nil
}

func (cn *conn) Rollback() error {
	return nil
}

func (cn *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	req, err := Parse(query)
	if err != nil {
		return nil, err
	}

	body, err := cn.cl.Post(req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	res, err := decode(req, body)
	if err != nil {
		return nil, err
	}

	return res, nil
}
