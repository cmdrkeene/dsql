package dsql

import (
	"database/sql"
	"io"
	"io/ioutil"

	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	name := "dyanmodb://access:secret@us-east-1"

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	Clients[name] = MockClient{
		OnPost: func(Operation, Request) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(`{"Count": 1}`)), nil
		},
	}

	rows, err := db.Query("SELECT id, email FROM users WHERE id=1;")
	if err != nil {
		t.Error(err)
	}

	if !rows.Next() {
		t.Fatal("expected row")
	}

	var id int
	var email string

	err = rows.Scan(&id, &email)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Error("bad id", id)
	}

	if email != "test@example.com" {
		t.Error("bad email", email)
	}
}
