package local

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ConfigName = "local"
)

type LocalConfig struct {
	ClientRepoConfig ClientRepositoryConfig `mapstructure:"clientRepo"`
	d                driver.Driver
	clientRepo       *ClientRepository
	httpMiddlewares  []func(http.Handler) http.Handler
}

func (l *LocalConfig) Init() (err error) {
	fmt.Printf("Initializing local configuration...")
	err = l.initClientRepo()
	if err != nil {
		return
	}
	return
}

func (l *LocalConfig) initClientRepo() (err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", l.ClientRepoConfig.dsn())
	if err != nil {
		return
	}

	l.clientRepo = NewClientRepository(db)
	return
}
