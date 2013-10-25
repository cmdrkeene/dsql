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

type result struct{}

func (r *result) LastInsertId() (int64, error) {
	// not very helpful in dynamo
	// request has this and data type can be string, etc
	return 0, nil
}

func (r *result) RowsAffected() (int64, error) {
	return 1, nil
}

func (cn *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	// do retry of transient errors here
	_, err := cn.query(query, args)
	if err != nil {
		return nil, err
	}

	return &result{}, nil
}

func (cn *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	// do retry of transient errors here
	return cn.query(query, args)
}

func (cn *conn) query(query string, args []driver.Value) (driver.Rows, error) {
	stmt := Statement{query, args}

	query = stmt.Prepare()

	req, err := Parse(query)
	if err != nil {
		return nil, err
	}

	body, err := cn.cl.Post(req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return req.Rows(body)
}
