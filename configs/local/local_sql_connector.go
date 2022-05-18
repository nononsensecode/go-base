package local

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
)

func (l *LocalConfig) SqlConnectorProvider() (pName string, p base.SqlConnectorProvider) {
	pName = "local"
	p = l
	return
}

func (l *LocalConfig) InitSqlDB(d driver.Driver) {
	var (
		err      error
		clientDb *sql.DB
	)

	clientDb, err = sql.Open("mysql", l.ClientRepoConfig.dsn())
	if err != nil {
		panic(err)
	}

	l.clientRepo = NewClientRepository(clientDb)
	l.d = d
}

type Connector struct {
	cId string
	r   *ClientRepository
	d   driver.Driver
	m   *sync.Mutex
}

func (c *Connector) Connect(ctx context.Context) (d driver.Conn, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	var dsn string

	dsn, err = c.r.GetDsnByClientId(c.cId)
	if err != nil {
		return
	}

	d, err = c.d.Open(dsn)
	return
}

func (c *Connector) Driver() driver.Driver {
	return c.d
}

func (l *LocalConfig) NewSqlConnector(ctx context.Context) (d driver.Connector, err error) {
	err = l.isInitialized()
	if err != nil {
		return
	}

	var (
		clientId string
		ok       bool
	)
	if clientId, ok = ctx.Value(ctxtypes.CtxClientIdKey).(string); !ok {
		err = fmt.Errorf("there is no client id in context")
		return
	}
	d = &Connector{
		cId: clientId,
		r:   l.clientRepo,
		d:   l.d,
		m:   &sync.Mutex{},
	}
	return
}
