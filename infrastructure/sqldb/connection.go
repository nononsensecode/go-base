package sqldb

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type SqlConnectionProvider interface {
	GetSqlConnection(ctx context.Context) (*sql.DB, error)
}

type DefaultSqlConnectionProvider struct{}

func (c DefaultSqlConnectionProvider) GetConnection(ctx context.Context) (db *sql.DB, err error) {
	var d driver.Connector
	d, err = GetSqlConnector(ctx)

	if err != nil {
		return
	}

	db = sql.OpenDB(d)

	if err = db.Ping(); err != nil {
		return
	}

	return
}
