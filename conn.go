// database/sql implementation
package dsql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"strconv"

	"regexp"
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
	if len(args) > 0 {
		prepared, err := statement(query).prepare(args)
		if err != nil {
			return nil, err
		}
		query = prepared
	}

	req, err := Parse(query)
	if err != nil {
		return nil, err
	}

	body, err := cn.cl.Post(req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	res, err := req.Result(body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type statement string

var statementPlaceholders = regexp.MustCompile("(\\$\\d+)+")

func (stmt statement) prepare(args []driver.Value) (prepared string, err error) {
	query := string(stmt)
	prepared = statementPlaceholders.ReplaceAllStringFunc(query, func(match string) string {
		offset, _ := strconv.Atoi(match[1:])
		return stmt.quote(args[offset-1])
	})
	return prepared, err
}

func (stmt statement) quote(v driver.Value) (quoted string) {
	switch v.(type) {
	case string, []byte:
		quoted = fmt.Sprintf("\"%v\"", v)
	default:
		quoted = fmt.Sprintf("%v", v)
	}
	return quoted
}
