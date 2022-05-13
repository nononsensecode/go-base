package mysqlconn

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
)

func GetConnection(ctx context.Context) (d *sqlx.DB, err error) {
	var conn driver.Connector
	conn, err = sqldb.GetConnector(ctx)
	if err != nil {
		return
	}

	d = sqlx.NewDb(sql.OpenDB(conn), "mysql")
	return
}
