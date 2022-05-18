package aws

import (
	"context"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nononsensecode/go-base"
)

func (a *AWSConfig) NewPgSqlPoolConnector(ctx context.Context) (connector base.PgSqlPoolConnector, err error) {
	err = a.isInitialized()
	if err != nil {
		return
	}

	var secretName string

	secretName, err = getSecretName(ctx)
	if err != nil {
		return
	}

	connector = &AWSPgSqlPoolConnector{
		c:          a.cache,
		secretName: secretName,
	}
	return
}

func (a *AWSConfig) PgSqlPoolConnectorProvider() (pName string, provider base.PgSqlPoolConnectorProvider) {
	pName = "aws"
	provider = a
	return
}

type AWSPgSqlPoolConnector struct {
	c          *secretcache.Cache
	secretName string
}

func (c *AWSPgSqlPoolConnector) GetPgSqlPool(ctx context.Context) (pool *pgxpool.Pool, err error) {
	var dbConfig DbConfig

	dbConfig, err = getDbConfig(c.c, c.secretName)
	if err != nil {
		return
	}

	pool, err = pgxpool.Connect(ctx, dbConfig.Dsn)
	return
}
