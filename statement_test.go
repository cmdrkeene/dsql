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
	actual, err := stmt.Prepare()
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Error("expected ", expected)
		t.Error("actual   ", actual)
	}
}

func BenchmarkStatementPrepareString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := &Statement{
			"INSERT INTO users (id, names) VALUES (?, ?)",
			[]driver.Value{1, "a"},
		}
		s.prepareString()
	}
}
