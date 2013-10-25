package dsql

import (
	"database/sql/driver"
	"testing"
)

func TestStatementPrepare(t *testing.T) {
	stmt := &Statement{
		"INSERT INTO users (id, names) VALUES (?, ?)",
		[]driver.Value{1, "a"},
	}

	expected := `INSERT INTO users (id, names) VALUES (1, "a")`
	actual := stmt.Prepare()

	if actual != expected {
		t.Error("expected ", expected)
		t.Error("actual   ", actual)
	}
}
