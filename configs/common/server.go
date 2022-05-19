package common

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
)

type ServerConfig struct {
	Persistence PersistenceConfig `mapstructure:"persistence"`
	Http        HttpConfig        `mapstructure:"http"`
	Log         LogConfig         `mapstructure:"log"`
}

func (s ServerConfig) Init() (err error) {
	if s.Http.isNil() {
		s.Http.Port = 8080
	}

	if s.Log.isNil() {
		s.Log.Level = "DEBUG"
	}

	if s.Persistence.SqlEnable && s.Persistence.Sql.isNil() {
		return fmt.Errorf("sql vendor is not specified")
	}

	if err = s.Persistence.Sql.init(); err != nil {
		return
	}

	return nil
}

func (h HttpConfig) isNil() bool {
	return h.Port == 0
}

func (h HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (s SqlConfig) SqlDbType() sqldb.DbType {
	return s.dbType
}

type PersistenceConfig struct {
	SqlEnable   bool      `mapstructure:"sqlEnable"`
	PgxEnable   bool      `mapstructure:"pgxEnable"`
	MongoEnable bool      `mapstructure:"mongoEnable"`
	Sql         SqlConfig `mapstructure:"sql"`
}

type SqlConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
	dbType    sqldb.DbType
	driver    driver.Driver
}

func (s *SqlConfig) isNil() bool {
	return s.SqlVendor == ""
}

func (s *SqlConfig) init() (err error) {
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

func (s SqlConfig) Driver() driver.Driver {
	return s.driver
}

func (s SqlConfig) ConnectionProvider() sqldb.SqlConnectionProvider {
	return SqlConnectionProviderImpl{
		dbType: s.dbType,
	}
}

type HttpConfig struct {
	Host               string   `mapstructure:"host"`
	Port               int      `mapstructure:"port"`
	ApiPrefix          string   `mapstructure:"apiPrefix"`
	CorsAllowedOrigins []string `mapstructure:"corsAllowedOrigins"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	IsDev bool   `mapstructure:"isDev"`
}

func (l LogConfig) isNil() bool {
	return l.Level == ""
}

type SqlConnectionProviderImpl struct {
	dbType sqldb.DbType
}

func (p SqlConnectionProviderImpl) GetSqlConnection(ctx context.Context) (db *sql.DB, err error) {
	var d driver.Connector
	d, err = sqldb.GetSqlConnector(ctx)
	if err != nil {
		return
	}

	db = sql.OpenDB(d)
	return
}
