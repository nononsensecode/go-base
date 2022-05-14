package sqldb

import (
	"context"
	"database/sql"
)

type ConnectionProvider interface {
	GetConnection(ctx context.Context) (*sql.DB, error)
}
