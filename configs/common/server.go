package common

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/nononsensecode/go-base/infrastructure/sqldb"
	"github.com/nononsensecode/go-base/logs"
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

	if err = s.Persistence.init(); err != nil {
		return
	}

	logMode := "production mode"
	if s.Log.IsDev {
		logMode = "development mode"
	}
	fmt.Printf("Initializing logger with default log level \"%s\" and logger is in \"%s\"\n", s.Log.Level, logMode)
	logs.Init(s.Log.Level, s.Log.IsDev)

	return nil
}

func (h HttpConfig) isNil() bool {
	return h.Port == 0
}

func (h HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
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
