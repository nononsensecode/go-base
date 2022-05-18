package base

import (
	"context"
	"database/sql/driver"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Configurer interface {
	Init()
}

type SqlConnectorProvider interface {
	NewSqlConnector(ctx context.Context) (driver.Connector, error)
	InitSqlDB(d driver.Driver)
	SqlConnectorProvider() (name string, provider SqlConnectorProvider)
}

type MiddlewareProvider interface {
	GetMiddlewares() []func(http.Handler) http.Handler
}

type PgSqlPoolConnectorProvider interface {
	NewPgSqlPoolConnector(ctx context.Context) (connector PgSqlPoolConnector, err error)
	PgSqlPoolConnectorProvider() (pName string, provider PgSqlPoolConnectorProvider)
}

type PgSqlPoolConnector interface {
	GetPgSqlPool(ctx context.Context) (pool *pgxpool.Pool, err error)
}
