package local

import (
	"database/sql/driver"
	"net/http"

	"github.com/jmoiron/sqlx"
	"gitlab.com/kaushikayanam/base/configs/common"
	"gitlab.com/kaushikayanam/base/infrastructure/sqldb"
)

type LocalConfig struct {
	ClientRepoConfig ClientRepositoryConfig `mapstructure:"clientRepo"`
	d                driver.Driver
	clientRepo       *ClientRepository
	httpMiddlewares  []func(http.Handler) http.Handler
}

func (l *LocalConfig) ConnectorProvider() (pName string, p sqldb.ConnectorProvider) {
	pName = "local"
	p = l
	return
}

func (l *LocalConfig) InitDB(dbName string) {
	var (
		err      error
		clientDb *sqlx.DB
	)

	clientDb, err = sqlx.Open("mysql", l.ClientRepoConfig.dsn())
	if err != nil {
		panic(err)
	}

	l.clientRepo = NewClientRepository(clientDb)

	l.d, err = common.InitDriver(dbName)
	if err != nil {
		panic(err)
	}

	l.httpMiddlewares = make([]func(http.Handler) http.Handler, 0)
}

func (l *LocalConfig) GetMiddlewares() []func(http.Handler) http.Handler {
	return l.httpMiddlewares
}
