package dsql

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Statement struct {
	query string
	args  []driver.Value
}

func (s *Statement) Prepare() (string, error) {
	// return s.prepareRegexp(), nil
	return s.prepareString(), nil
}

func (s *Statement) prepareString() string {
	for _, v := range s.args {
		s.query = strings.Replace(s.query, "?", s.quote(v), 1)
	}
	return s.query
}

func (s *Statement) quote(v driver.Value) string {
	switch v.(type) {
	case string, []byte:
		return fmt.Sprintf("\"%v\"", v)
	}
	return fmt.Sprintf("%v", v)
}
