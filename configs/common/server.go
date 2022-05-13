package common

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
)

type ServerConfig struct {
	Persistence PersistenceConfig `mapstructure:"persistence"`
	Http        HttpConfig        `mapstructure:"http"`
	Log         LogConfig         `mapstructure:"log"`
}

func (s ServerConfig) Init() error {
	return s.Persistence.init()
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Http.Host, s.Http.Port)
}

func (pc PersistenceConfig) SqlDbType() sqldb.DbType {
	return pc.dbType
}

func (pc PersistenceConfig) SqlDriver() driver.Driver {
	return pc.driver
}

type PersistenceConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
	dbType    sqldb.DbType
	driver    driver.Driver
}

func (pc *PersistenceConfig) init() (err error) {
	pc.dbType, err = sqldb.NewDbType(pc.SqlVendor)
	if err != nil {
		return
	}

	pc.driver, err = pc.initSqlDriver()
	if err != nil {
		return
	}

	return
}

func (pc PersistenceConfig) initSqlDriver() (d driver.Driver, err error) {
	switch pc.dbType.String() {
	case "mysql":
		d = mysql.MySQLDriver{}
		return
	case "sqlite":
		d = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named %s", pc.dbType.String())
	}
	return
}

func (pc PersistenceConfig) GetConnection(ctx context.Context) (db *sqlx.DB, err error) {
	var d driver.Connector
	d, err = sqldb.GetConnector(ctx)

	if err != nil {
		return
	}

	db = sqlx.NewDb(sql.OpenDB(d), pc.dbType.String())
	return
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
