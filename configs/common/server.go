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

func (h HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (pc PersistenceConfig) SqlDbType() sqldb.DbType {
	return pc.dbType
}

type PersistenceConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
	dbType    sqldb.DbType
}

func (pc *PersistenceConfig) SqlInit() (d driver.Driver, err error) {
	pc.dbType, err = sqldb.NewDbType(pc.SqlVendor)
	if err != nil {
		return
	}

	switch pc.dbType.String() {
	case "mysql":
		fmt.Println("configured sql driver is \"mysql\"")
		d = mysql.MySQLDriver{}
		return
	case "sqllite":
		fmt.Printf("configured sql driver is \"sqllite\"")
		d = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named \"%s\"", pc.dbType.String())
	}
	return
}

func (pc PersistenceConfig) ConnectionProvider() sqldb.ConnectionProvider {
	return ConnectionProviderImpl{
		dbType: pc.dbType,
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

type ConnectionProviderImpl struct {
	dbType sqldb.DbType
}

func (p ConnectionProviderImpl) GetConnection(ctx context.Context) (db *sql.DB, err error) {
	var d driver.Connector
	d, err = sqldb.GetConnector(ctx)
	if err != nil {
		return
	}

	db = sql.OpenDB(d)
	return
}
