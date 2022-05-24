package common

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
)

type PersistenceConfig struct {
	SqlEnable   bool      `mapstructure:"sqlEnable"`
	PgxEnable   bool      `mapstructure:"pgxEnable"`
	MongoEnable bool      `mapstructure:"mongoEnable"`
	Sql         SqlConfig `mapstructure:"sql"`
}

func (p PersistenceConfig) init() (err error) {
	if p.SqlEnable && p.PgxEnable {
		err = fmt.Errorf("at a time one sql vendor can be enabled")
		return
	}

	if p.SqlEnable && p.Sql.isNil() {
		err = fmt.Errorf("sql configuration is empty")
		return
	}

	err = p.Sql.init()
	if err != nil {
		return
	}
	return
}

type SqlConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
	dbType    sqldb.DbType
	driver    driver.Driver
}

func (s SqlConfig) SqlDbType() sqldb.DbType {
	return s.dbType
}

func (s *SqlConfig) isNil() bool {
	return s.SqlVendor == ""
}

func (s *SqlConfig) init() (err error) {
	fmt.Println("Initializing sql....")
	s.dbType, err = sqldb.NewDbType(s.SqlVendor)
	if err != nil {
		return
	}

	switch s.dbType.String() {
	case "mysql":
		fmt.Println("configured sql driver is \"mysql\"")
		s.driver = mysql.MySQLDriver{}
		return
	case "sqllite":
		fmt.Printf("configured sql driver is \"sqllite\"")
		s.driver = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named \"%s\"", s.dbType.String())
	}
	return
}

func (s *SqlConfig) Driver() driver.Driver {
	return s.driver
}

func (s *SqlConfig) ConnectionProvider() sqldb.SqlConnectionProvider {
	return SqlConnectionProviderImpl{
		dbType: s.dbType,
	}
}
