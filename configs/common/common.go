package common

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

func InitDriver(name string) (d driver.Driver, err error) {
	switch name {
	case "mysql":
		d = mysql.MySQLDriver{}
		return
	case "sqlite":
		d = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named %s", name)
	}
	return
}
