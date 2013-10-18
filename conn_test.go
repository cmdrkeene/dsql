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

	rows, err := db.Query("SELECT id, name FROM users WHERE id=1;")
	if err != nil {
		t.Error(err)
	}

	t.Error(rows)
}
