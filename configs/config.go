package configs

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/kaushikayanam/base/configs/aws"
	"gitlab.com/kaushikayanam/base/configs/common"
	"gitlab.com/kaushikayanam/base/configs/local"
	"gitlab.com/kaushikayanam/base/infrastructure/sqldb"
	"gitlab.com/kaushikayanam/base/interfaces/httpsrvr"
)

type Config struct {
	Server                  common.ServerConfig    `mapstructure:"server"`
	Configs                 map[string]interface{} `mapstructure:"configs"`
	isInitialized           bool
	httpMiddlewareProviders []httpsrvr.MiddlewareProvider
}

func (cfg *Config) Init() {
	cfg.isInitialized = true

	var connProviders []sqldb.ConnectorProvider

	for cName, configuration := range cfg.Configs {
		switch cName {
		case "local":
			var lcfg local.LocalConfig
			if err := mapstructure.Decode(configuration, &lcfg); err != nil {
				panic(err)
			}
			lcfg.InitDB(cfg.Server.Persistence.Vendor)
			connProviders = append(connProviders, &lcfg)
			cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &lcfg)
		case "aws":
			var acfg aws.AWSConfig
			if err := mapstructure.Decode(configuration, &acfg); err != nil {
				panic(err)
			}
			acfg.InitDB(cfg.Server.Persistence.Vendor)
			connProviders = append(connProviders, &acfg)
			// cfg.httpMiddlewareProviders = append(cfg.httpMiddlewareProviders, &acfg)
		default:
			panic(fmt.Errorf("invalid configuration %s", cName))
		}
	}

	sqldb.Init(connProviders)
}
