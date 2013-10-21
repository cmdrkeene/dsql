// Declare types and conversion maps
package dsql

import "strconv"

type Operation string
type Request interface{}

var TokenTypes = map[Token]string{
	Number: "N",
	String: "S",
}

var DefinitionTypes = map[string]string{
	"string":    "S",
	"stringset": "SS",
	"number":    "N",
	"numberset": "NS",
	"binary":    "B",
	"binaryset": "BS",
}

func NewValue(t Token, s string) Value {
	return Value{TokenTypes[t], trim(s)}
}

type Value struct {
	Type  string
	Value string
}

type Item struct {
	S string `json:"S,omitempty`
	N string `json:"N,omitempty`
}

func (i Item) Value() interface{} {
	if i.S != "" {
		return []byte(i.S)
	}

	if i.N != "" {
		i, _ := strconv.Atoi(i.N)
		return i
	}

	return nil
}
