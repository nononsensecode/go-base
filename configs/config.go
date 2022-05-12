package configs

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/configs/aws"
	"github.com/nononsensecode/go-base/configs/common"
	"github.com/nononsensecode/go-base/configs/local"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
	"github.com/nononsensecode/go-base/interfaces/httpsrvr"
	"github.com/nononsensecode/go-base/logs"
)

type Config struct {
	Server                  common.ServerConfig `mapstructure:"server"`
	PlatformConfig          PlatformConfig      `mapstructure:"platform"`
	httpMiddlewareProviders []base.MiddlewareProvider
	isInitialized           bool
}

func (cfg *Config) Init() {
	cfg.isInitialized = true
	httpsrvr.Middlewares = cfg.getHttpMiddlewares()
	cfg.initsql()
	cfg.InitLogger()
}

func (cfg *Config) InitSqlDriver() (d driver.Driver, err error) {
	switch cfg.Server.SqlVendor() {
	case "mysql":
		d = mysql.MySQLDriver{}
		return
	case "sqlite":
		d = &sqlite3.SQLiteDriver{}
	default:
		err = fmt.Errorf("there is no sql driver named %s", cfg.Server.SqlVendor())
	}
	return
}

func (cfg *Config) initsql() {
	var (
		connProviders []base.ConnectorProvider
		d             driver.Driver
		err           error
	)

	if d, err = cfg.InitSqlDriver(); err != nil {
		panic(err)
	}

	switch cfg.PlatformConfig.Name {
	case "local":
		var lcfg local.LocalConfig
		if err := mapstructure.Decode(cfg.PlatformConfig.Config, &lcfg); err != nil {
			panic(err)
		}
		lcfg.InitDB(d)
		connProviders = append(connProviders, &lcfg)
		cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &lcfg)
	case "aws":
		var acfg aws.AWSConfig
		if err := mapstructure.Decode(cfg.PlatformConfig.Config, &acfg); err != nil {
			panic(err)
		}
		acfg.InitDB(d)
		connProviders = append(connProviders, &acfg)
		cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &acfg)
	default:
		panic(fmt.Errorf("invalid configuration %s", cfg.PlatformConfig.Name))
	}

	sqldb.Init(connProviders)
}

func (cfg *Config) InitLogger() {
	logs.Init(cfg.Server.Log.Level, cfg.Server.Log.IsDev)
}

type PlatformConfig struct {
	Name   string      `mapstructure:"name"`
	Config interface{} `mapstructure:"config"`
}
