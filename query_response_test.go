package dsql

import (
	"reflect"
	"testing"
)

func TestQueryResponseColumns(t *testing.T) {
	r := &QueryResponse{
		Count: 3,
		Items: []Item{
			Item{
				"id":    Attribute{N: "1"},
				"email": Attribute{S: "a@example.com"},
			},
			Item{
				"id":    Attribute{N: "1"},
				"email": Attribute{S: "a@example.com"},
				"bio":   Attribute{S: "person"},
			},
			Item{
				"height": Attribute{S: "7'4"},
				"weight": Attribute{S: "520 LB'"},
			},
		},
	}

	expected := []string{"id", "email", "bio", "height", "weight"}
	actual := r.Columns()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("expected ", expected)
		t.Error("actual   ", actual)
	}
}
