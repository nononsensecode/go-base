package sqldb

import (
	"database/sql"
	"fmt"

	"go.uber.org/multierr"
)

func FinishTransaction(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			rollbackErr = fmt.Errorf("rolling back failed: %w", rollbackErr)
			return multierr.Combine(rollbackErr, err)
		}

		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		commitErr = fmt.Errorf("committing failed: %w", commitErr)
		return multierr.Combine(commitErr, err)
	}

	return nil
}
