package common

import (
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

func (s ServerConfig) Init() error {
	return s.Persistence.init()
}

func (h HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
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

	if err = pc.initSqlDriver(); err != nil {
		return
	}

	return
}

func (pc *PersistenceConfig) initSqlDriver() (err error) {
	switch pc.dbType.String() {
	case "mysql":
		fmt.Println("configured sql driver is \"mysql\"")
		pc.driver = mysql.MySQLDriver{}
		return
	case "sqllite":
		fmt.Printf("configured sql driver is \"sqllite\"")
		pc.driver = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named \"%s\"", pc.dbType.String())
	}

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
