package sqldb

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ConnectionProvider interface {
	GetConnection(ctx context.Context) (*sqlx.DB, error)
}
