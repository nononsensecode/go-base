package common

import (
	"fmt"

	"github.com/go-chi/cors"
)

type ServerConfig struct {
	Host        string            `mapstructure:"host"`
	Port        int               `mapstructure:"port"`
	ApiPrefix   string            `mapstructure:"apiPrefix"`
	Persistence PersistenceConfig `mapstructure:"persistence"`
	Cors        CorsConfig        `mapstructure:"cors"`
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s ServerConfig) SqlVendor() string {
	return s.Persistence.SqlVendor
}

func (s ServerConfig) ChiCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:   s.Cors.AllowedAddresses,
		AllowedMethods:   s.Cors.AllowedMethods,
		AllowedHeaders:   s.Cors.AllowedHeaders,
		ExposedHeaders:   s.Cors.ExposedHeaders,
		AllowCredentials: s.Cors.AllowCredentials,
		MaxAge:           s.Cors.MaxAge,
	}
}

type CorsConfig struct {
	AllowedAddresses []string `mapstructure:"allowedAddresses"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposedHeaders   []string `mapstructure:"exposedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

type PersistenceConfig struct {
	SqlVendor string `mapstructure:"sqlVendor"`
}
