package base

import (
	"context"
	"database/sql/driver"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type Configurer interface {
	Init() error
}

type MiddlewareProvider interface {
	GetMiddlewares() []func(http.Handler) http.Handler
}

type SqlConnectorProvider interface {
	NewSqlConnector(ctx context.Context) (driver.Connector, error)
	InitSqlDB(d driver.Driver)
	SqlConnectorProvider() (name string, provider SqlConnectorProvider)
}

type PgSqlPoolConnectorProvider interface {
	NewPgSqlPoolConnector(ctx context.Context) (connector PgSqlPoolConnector, err error)
	PgSqlPoolConnectorProvider() (pName string, provider PgSqlPoolConnectorProvider)
}

type PgSqlPoolConnector interface {
	GetPgSqlPool(ctx context.Context) (pool *pgxpool.Pool, err error)
}

type MongoDbClientBuilderProvider interface {
	NewMongoDbClientBuilder(ctx context.Context) (builder MongoDbClientBuilder, err error)
	MongoDbClientBuilderProvider() (pName string, provider MongoDbClientBuilderProvider)
}

type MongoDbClientBuilder interface {
	GetMongoDbClient(ctx context.Context) (client *mongo.Client, err error)
}
