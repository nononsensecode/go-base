package local

import (
	"database/sql/driver"
	"fmt"
	"net/http"
)

type LocalConfig struct {
	ClientRepoConfig ClientRepositoryConfig `mapstructure:"clientRepo"`
	d                driver.Driver
	clientRepo       *ClientRepository
	httpMiddlewares  []func(http.Handler) http.Handler
}

func (l *LocalConfig) isInitialized() (err error) {
	if l.clientRepo == nil {
		err = fmt.Errorf("local configuration is not initialized")
	}
	return
}
