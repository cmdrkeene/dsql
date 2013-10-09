package sql

type Select struct{}
type Insert struct{}
type Update struct{}
type Delete struct{}
type Create struct{}
type Drop struct{}

func Translate(sql string) (json string, err error) {
	return "", nil
}
