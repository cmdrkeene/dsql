package dsql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestQuerySelect(t *testing.T) {
	name := "dyanmodb://access:secret@us-east-1"

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	Clients[name] = MockClient{
		OnPost: func(Request) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(`{
				"Count": 1,
				"Items": [
					{
						"id": {"N": "1"},
						"email": {"S": "test@example.com"}
					}
				]
			}`)), nil
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
		t.Error("id", id)
	}

	if email != "test@example.com" {
		t.Error("email", email)
	}
}

func TestQueryClientError(t *testing.T) {
	name := "dyanmodb://access:secret@us-east-1"

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	e := errors.New("something went wrong")

	Clients[name] = MockClient{
		OnPost: func(Request) (io.ReadCloser, error) {
			return nil, e
		},
	}

	_, err = db.Query(`INSERT INTO users (id, name) VALUES (1, "a")`)

	if err != e {
		t.Error("expected ", e)
		t.Error("actual   ", err)
	}
}

func TestQueryInsert(t *testing.T) {
	name := "dyanmodb://access:secret@us-east-1"

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	Clients[name] = MockClient{
		OnPost: func(Request) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(`{}`)), nil
		},
	}

	_, err = db.Query(`INSERT INTO users (id, name) VALUES (1, "a")`)
	if err != nil {
		t.Error(err)
	}
}

func TestStatementPrepare(t *testing.T) {
	s := statement("INSERT INTO users (id, names) VALUES ($1, $2)")

	expected := `INSERT INTO users (id, names) VALUES (1, "a")`
	actual, err := s.prepare([]driver.Value{1, "a"})
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("expected ", expected)
		t.Error("actual   ", actual)
	}
}

type MockClient struct {
	OnPost func(Request) (io.ReadCloser, error)
}

func (m MockClient) Post(r Request) (io.ReadCloser, error) {
	return m.OnPost(r)
}
