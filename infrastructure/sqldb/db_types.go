package sqldb

import "fmt"

type DbType int

const (
	MySql DbType = iota
	SqlLite
	Postgresql
)

var d = [...]string{"mysql", "sqllite", "postgresql"}

func (s DbType) String() string {
	return d[s]
}

func NewDbType(t string) (dt DbType, err error) {
	for i, v := range d {
		if v == t {
			return DbType(i), nil
		}
	}

	return -1, fmt.Errorf("there is no db type named \"%s\"", t)
}
