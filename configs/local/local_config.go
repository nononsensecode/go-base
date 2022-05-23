package local

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
	"github.com/nononsensecode/go-base/interfaces/httpsrvr"
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

func (l *LocalConfig) Init(sqlDriver driver.Driver) (err error) {
	fmt.Printf("Initializing local configuration...")
	if sqlDriver == nil {
		err = fmt.Errorf("sql driver is nil")
		return
	}

	err = l.initClientRepo()
	if err != nil {
		return
	}

	l.d = sqlDriver

	sqldb.Init(l)
	httpsrvr.AddMiddlewares(l.GetMiddlewares()...)

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
