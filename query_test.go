package dsql

import (
	"database/sql/driver"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestQueryResult(t *testing.T) {
	r := ioutil.NopCloser(strings.NewReader(`
  {
    "Count": 1,
    "Items": [
      {
        "id": {"N": "1"},
        "email": {"S": "test@example.com"},
        "bio": null
      }
    ]
  }`))

	q := Query{}
	res, err := q.Result(r)
	if err != nil {
		t.Error(err)
	}

	cols := []string{"id", "email", "bio"}
	if !reflect.DeepEqual(res.Columns(), cols) {
		t.Error("bad columns", res.Columns())
	}

	dest := make([]driver.Value, len(cols))
	err = res.Next(dest)
	if err != nil {
		t.Error(err)
	}

	expected := []driver.Value{1, []byte("test@example.com"), nil}
	if !reflect.DeepEqual(dest, expected) {
		t.Error("expected ", expected)
		t.Error("actual   ", dest)
	}

	err = res.Next(dest)
	if err != io.EOF {
		t.Error("expected EOF", err, dest)
	}
}
