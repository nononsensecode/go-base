package common

import (
	"fmt"
)

type ServerConfig struct {
	Host        string            `mapstructure:"host"`
	Port        int               `mapstructure:"port"`
	ApiPrefix   string            `mapstructure:"apiPrefix"`
	Persistence PersistenceConfig `mapstructure:"persistence"`
	Http        HttpConfig        `mapstructure:"http"`
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s ServerConfig) SqlVendor() string {
	return s.Persistence.SqlVendor
}

type PersistenceConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
}

type HttpConfig struct {
	CorsAllowedOrigins []string `mapstructure:"corsAllowedOrigins"`
}
