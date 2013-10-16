package dsql

import "strings"

func trim(s string) string {
	return strings.Trim(s, `"`)
}
