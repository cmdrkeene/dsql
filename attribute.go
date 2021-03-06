package dsql

import "strconv"

type Attribute struct {
	S string `json:",omitempty"`
	N string `json:",omitempty"`
}

func (a Attribute) Value() interface{} {
	if a.S != "" {
		return []byte(a.S)
	}

	if a.N != "" {
		i, _ := strconv.Atoi(a.N)
		return i
	}

	return nil
}
