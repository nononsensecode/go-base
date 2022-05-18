package sqldb

import (
	"context"
	"database/sql"
)

type SqlConnectionProvider interface {
	GetSqlConnection(ctx context.Context) (*sql.DB, error)
}
