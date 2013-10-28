package dsql

import (
	"database/sql"

	"errors"
	"io"
	"io/ioutil"

	"reflect"
	"strings"
	"testing"
)

func TestExecInsert(t *testing.T) {
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

	result, err := db.Exec("INSERT INTO users (id, email) VALUES (?, ?)", 1, "test@example.com")
	if err != nil {
		t.Error(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Error(err)
	}
	if id != 0 {
		t.Error("expected ", 0)
		t.Error("actual   ", id)
	}

	n, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Error("expected ", 1)
		t.Error("actual   ", n)
	}
}

func TestQuerySelect(t *testing.T) {
	name := "dyanmodb://access:secret@us-east-1"

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	Clients[name] = MockClient{
		OnPost: func(Request) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(`{
				"Count": 2,
				"Items": [
					{
						"id": {"N": "1"},
						"email": {"S": "test@example.com"}
					},
					{
						"id": {"N": "2"}
					}
				]
			}`)), nil
		},
	}

	rows, err := db.Query("SELECT id, email FROM users WHERE id=1;")
	if err != nil {
		t.Error(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(cols, []string{"id", "email"}) {
		t.Error("bad columns", cols)
	}

	// row 1

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

	// row 2

	if !rows.Next() {
		t.Fatal("expected row")
	}

	var id2 int
	var email2 string

	err = rows.Scan(&id2, &email2)
	if err != nil {
		t.Fatal(err)
	}

	if id2 != 2 {
		t.Error("id", id2)
	}

	if email2 != "" {
		t.Error("email", email2)
	}

	// eof
	if rows.Next() {
		t.Error("should be EOF")
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

type MockClient struct {
	OnPost func(Request) (io.ReadCloser, error)
}

func (m MockClient) Post(r Request) (io.ReadCloser, error) {
	return m.OnPost(r)
}
