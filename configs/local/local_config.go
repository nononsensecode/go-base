package local

import (
	"database/sql/driver"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nononsensecode/go-base"
)

type LocalConfig struct {
	ClientRepoConfig ClientRepositoryConfig `mapstructure:"clientRepo"`
	d                driver.Driver
	clientRepo       *ClientRepository
	httpMiddlewares  []func(http.Handler) http.Handler
}

func (l *LocalConfig) ConnectorProvider() (pName string, p base.ConnectorProvider) {
	pName = "local"
	p = l
	return
}

func (l *LocalConfig) InitDB(d driver.Driver) {
	var (
		err      error
		clientDb *sqlx.DB
	)

	clientDb, err = sqlx.Open("mysql", l.ClientRepoConfig.dsn())
	if err != nil {
		panic(err)
	}

	l.clientRepo = NewClientRepository(clientDb)
	l.d = d
}

// For future use
func (l *LocalConfig) GetMiddlewares() []func(http.Handler) http.Handler {
	return l.httpMiddlewares
}
