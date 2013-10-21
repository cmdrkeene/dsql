package dsql

import "strings"

func trim(s string) string {
	return strings.Trim(s, `"`)
}

func operation(r Request) Operation {
	switch r.(type) {
	case Query:
		return Operation("DynamoDB_20120810.Query")
	}
	return ""
}
