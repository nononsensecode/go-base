package common

import (
	"fmt"

	"github.com/nononsensecode/go-base"
)

type PlatformConfig struct {
	Name   string      `mapstructure:"name"`
	Config interface{} `mapstructure:"config"`
}

func (p PlatformConfig) isNil() bool {
	return p.Name == "" || p.Config == nil
}

func (p PlatformConfig) Init() (err error) {
	if p.isNil() {
		return fmt.Errorf("local/cloud configuration is missing")
	}

	var (
		c  base.Configurer
		ok bool
	)

	if c, ok = p.Config.(base.Configurer); !ok {
		return fmt.Errorf("configuration is missing initialization method")
	}

	err = c.Init()
	return
}
