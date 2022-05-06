package base

import (
	"context"
	"database/sql/driver"
	"net/http"
)

type ConnectorProvider interface {
	NewConnector(ctx context.Context) (driver.Connector, error)
	InitDB(d driver.Driver)
	ConnectorProvider() (string, ConnectorProvider)
}

type MiddlewareProvider interface {
	GetMiddlewares() []func(http.Handler) http.Handler
}
