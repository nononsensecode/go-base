package configs

import (
	"database/sql/driver"
	"fmt"

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
	Server                  common.ServerConfig   `mapstructure:"server"`
	Platform                common.PlatformConfig `mapstructure:"platform"`
	httpMiddlewareProviders []base.MiddlewareProvider
	isInitialized           bool
}

func (cfg *Config) Init() {
	cfg.isInitialized = true
	httpsrvr.Middlewares = cfg.getHttpMiddlewares()
	cfg.initsql()
	cfg.initLogger()
}

func (cfg *Config) initsql() {
	var (
		connProviders []base.SqlConnectorProvider
		d             driver.Driver
		err           error
	)

	d, err = cfg.Server.Persistence.SqlInit()
	if err != nil {
		panic(err)
	}

	switch cfg.PlatformConfig.Name {
	case "local":
		var lcfg local.LocalConfig
		if err := mapstructure.Decode(cfg.PlatformConfig.Config, &lcfg); err != nil {
			panic(err)
		}
		lcfg.InitSqlDB(d)
		connProviders = append(connProviders, &lcfg)
		cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &lcfg)
	case "aws":
		var acfg aws.AWSConfig
		if err := mapstructure.Decode(cfg.PlatformConfig.Config, &acfg); err != nil {
			panic(err)
		}
		acfg.InitSqlDB(d)
		connProviders = append(connProviders, &acfg)
		cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &acfg)
	default:
		panic(fmt.Errorf("invalid configuration %s", cfg.PlatformConfig.Name))
	}

	sqldb.Init(connProviders)
}

func (cfg *Config) initLogger() {
	logs.Init(cfg.Server.Log.Level, cfg.Server.Log.IsDev)
}
