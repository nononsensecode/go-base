package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
)

type LocalPgSqlPoolConnector struct {
	r        *ClientRepository
	m        *sync.Mutex
	clientId string
}

func (c *LocalPgSqlPoolConnector) GetPgSqlPool(ctx context.Context) (pool *pgxpool.Pool, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	var dsn string
	dsn, err = c.r.GetDsnByClientId(c.clientId)
	if err != nil {
		return
	}

	pool, err = pgxpool.Connect(ctx, dsn)
	return
}

func (l *LocalConfig) NewPgSqlPoolConnector(ctx context.Context) (connector base.PgSqlPoolConnector, err error) {
	if err = l.isInitialized(); err != nil {
		return
	}

	var (
		clientId string
		ok       bool
	)

	if clientId, ok = ctx.Value(ctxtypes.CtxClientIdKey).(string); !ok {
		err = fmt.Errorf("client id is missing in the context")
		return
	}

	connector = &LocalPgSqlPoolConnector{
		m:        &sync.Mutex{},
		r:        l.clientRepo,
		clientId: clientId,
	}
	return
}

func (l *LocalConfig) PgSqlPoolConnectorProvider() (pName string, provider base.PgSqlPoolConnectorProvider) {
	pName = "local"
	provider = l
	return
}
