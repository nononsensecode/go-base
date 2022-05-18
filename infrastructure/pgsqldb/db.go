package pgsqldb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
)

func FinishTransaction(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			rollBackErr = fmt.Errorf("rolling back failed: %w", rollBackErr)
			err = multierr.Append(err, rollBackErr)
		}
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		err = fmt.Errorf("committing failed: %w", commitErr)
	}
	return err
}
