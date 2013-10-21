/*
  Package dsql is a SQL dialect for interacting with Amazon's DynamoDB.

  Many methods are no-ops due to the lack of transactions, statements, etc.

  Example:

    import (
      _ "github.com/cmdrkeene/dsql"
      "database/sql"
    )

    func main() {
      url := "dynamodb://access_key:secret_key@us-east-1"
      db, _ := sql.Open("dynamodb", url)
      rows, _ := db.Query("SELECT name FROM users WHERE id=$1", 123)
    }
*/
package dsql

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log"
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

	body, err := cn.cl.Post(operation(req), req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = req.Decode(body)
	if err != nil {
		return nil, err
	}

	return &rows{req: req}, nil
}

type rows struct {
	req  Request
	done bool
}

func (r *rows) Columns() []string {
	return []string{"id", "email"}
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.done {
		log.Print("done")
		return io.EOF
	} else {
		dest[0] = 1
		dest[1] = "test@example.com"
		r.done = true
	}

	return nil
}
