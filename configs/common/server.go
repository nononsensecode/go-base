package common

import (
	"fmt"
)

type ServerConfig struct {
	Persistence PersistenceConfig `mapstructure:"persistence"`
	Http        HttpConfig        `mapstructure:"http"`
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Http.Host, s.Http.Port)
}

func (s ServerConfig) SqlVendor() string {
	return s.Persistence.SqlVendor
}

type PersistenceConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
}

type HttpConfig struct {
	Host               string   `mapstructure:"host"`
	Port               int      `mapstructure:"port"`
	ApiPrefix          string   `mapstructure:"apiPrefix"`
	CorsAllowedOrigins []string `mapstructure:"corsAllowedOrigins"`
}
