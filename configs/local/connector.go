package local

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/nononsensecode/go-base/context/ctxtypes"
)

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

func (l *LocalConfig) NewConnector(ctx context.Context) (d driver.Connector, err error) {
	if l.clientRepo == nil {
		err = fmt.Errorf("client repo is not initialized")
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
