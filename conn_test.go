package dsql

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("response")
		fmt.Fprint(w, `
      {
        "ConsumedCapacity": {"CapacityUnits": 1,"TableName": "users"}
        "Item": {
          "id": {"N": "123"},
          "name": {"S": "Brandon"},
        }
      }
    `)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	name := fmt.Sprintf("dyanmodb://access:secret@%s", u.Host)

	db, err := sql.Open("dynamodb", name)
	if err != nil {
		t.Error(err)
	}

	rows, err := db.Query("SELECT id, name FROM users WHERE id=1;")
	if err != nil {
		t.Error(err)
	}

	t.Error(rows)
}
