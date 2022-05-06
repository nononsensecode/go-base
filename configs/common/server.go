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

type CorsConfig struct {
	AllowedAddresses []string `mapstructure:"allowedAddresses"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposedHeaders   []string `mapstructure:"exposedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

func (c CorsConfig) ChiCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:   c.AllowedAddresses,
		AllowedMethods:   c.AllowedMethods,
		AllowedHeaders:   c.AllowedHeaders,
		ExposedHeaders:   c.ExposedHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	}
}

type PersistenceConfig struct {
	Vendor string `mapstructure:"vendor"`
}
