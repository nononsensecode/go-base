package pgsqldb

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nononsensecode/go-base"
)

type PgSqlConnectionPoolProvider interface {
	GetPgSqlConnectionPool(ctx context.Context) (*pgxpool.Pool, error)
}

type DefaultPgSqlConnectionPoolProvider struct{}

func (p DefaultPgSqlConnectionPoolProvider) GetPgSqlConnectionPool(ctx context.Context) (pool *pgxpool.Pool, err error) {
	var connector base.PgSqlPoolConnector
	connector, err = GetPgSqlConnectionPoolConnector(ctx)
	if err != nil {
		return
	}

	pool, err = connector.GetPgSqlPool(ctx)
	return
}
