package dsql

import (
	"encoding/json"
	"reflect"

	"testing"
)

func TestAttributeMarshal(t *testing.T) {
	s := Attribute{S: "a", N: ""}
	b, _ := json.Marshal(s)
	if !reflect.DeepEqual(b, []byte(`{"S":"a"}`)) {
		t.Error("bad marshal", string(b))
	}

	n := Attribute{S: "", N: "1"}
	b, _ = json.Marshal(n)
	if !reflect.DeepEqual(b, []byte(`{"N":"1"}`)) {
		t.Error("bad marshal", string(b))
	}
}
