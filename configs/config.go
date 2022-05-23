package configs

import (
	"github.com/nononsensecode/go-base/configs/common"
	"github.com/nononsensecode/go-base/interfaces/httpsrvr"
)

type Config struct {
	Server        common.ServerConfig   `mapstructure:"server"`
	Platform      common.PlatformConfig `mapstructure:"platform"`
	isInitialized bool
}

func (cfg *Config) Init() {
	if err := cfg.Server.Init(); err != nil {
		panic(err)
	}

	httpsrvr.AddMiddlewares(cfg.getHttpMiddlewares()...)

	d := cfg.Server.Persistence.Sql.Driver()
	if err := cfg.Platform.Init(d); err != nil {
		panic(err)
	}
}
