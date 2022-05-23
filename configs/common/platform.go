package common

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/configs/aws"
	"github.com/nononsensecode/go-base/configs/local"
)

type PlatformConfig struct {
	Name       string      `mapstructure:"name"`
	Config     interface{} `mapstructure:"config"`
	configurer base.Configurer
}

func (p PlatformConfig) isNil() bool {
	return p.Name == "" || p.Config == nil
}

func (p PlatformConfig) Init() (err error) {
	fmt.Println("Finding correct platform....")
	if p.isNil() {
		return fmt.Errorf("local/cloud configuration is missing")
	}
	fmt.Printf("Platform \"%s\" is getting unmarshalled....\n", p.Name)

	if p.configurer, err = p.decode(); err != nil {
		return
	}

	fmt.Printf("Initializing platform \"%s\"....\n", p.Name)
	err = p.configurer.Init()
	return
}

func (p PlatformConfig) GetConfigurer() base.Configurer {
	return p.configurer
}

func (p *PlatformConfig) decode() (c base.Configurer, err error) {
	switch p.Name {
	case "local":
		var l local.LocalConfig
		err = mapstructure.Decode(p.Config, &l)
		if err != nil {
			return
		}
		c = &l
		return
	case "aws":
		var a aws.AWSConfig
		err = mapstructure.Decode(p.Config, &a)
		if err != nil {
			return
		}
		c = &a
		return
	default:
		err = fmt.Errorf("there is no configuration named \"%s\"", p.Name)
	}
	return
}
