/*
  Package dsql is a SQL dialect for interacting with Amazon's DynamoDB.

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
	"github.com/groupme/dynamo"
	"log"
	"net/url"
	"strings"
)

func init() {
	sql.Register("dynamodb", &drv{})
}

type drv struct{}

func (d *drv) Open(name string) (driver.Conn, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	accessKey := u.User.Username()
	secretKey, _ := u.User.Password()
	db := dynamo.Open(accessKey, secretKey, u.Host)

	// for use with net/http/httptest
	if strings.HasPrefix(u.Host, "localhost") {
		db.URL = u.Host
	}

	return &conn{db: db}, nil
}

type conn struct {
	db *dynamo.DB
}

func (cn *conn) Prepare(q string) (driver.Stmt, error) {
	log.Print("prepare")
	return nil, nil
}

func (cn *conn) Close() error {
	return nil
}

func (cn *conn) Begin() (driver.Tx, error) {
	log.Print("begin")
	return cn, nil
}

func (cn *conn) Commit() error {
	log.Print("commit")
	return nil
}

func (cn *conn) Rollback() error {
	log.Print("rollback")
	return nil
}

// unsure if this is needed
// could be useful where you just don't care about the result
// e.g. DeleteTable, DeleteItem
// I wonder if that Exec/close problem is relevant here
func (cn *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (cn *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	// encode args, if any
	// parse query
	req, err := Parse(query)
	if err != nil {
		return nil, err
	}
	log.Printf("%#v", req)

	// make request
	// convert to rows
	return nil, nil
}
