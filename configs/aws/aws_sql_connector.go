package aws

import (
	"context"
	"database/sql/driver"
	"sync"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/nononsensecode/go-base"
)

type Connector struct {
	secretName string
	cache      *secretcache.Cache
	d          driver.Driver
	m          *sync.Mutex
}

func (c *Connector) Connect(ctx context.Context) (conn driver.Conn, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	var dbConfig DbConfig
	dbConfig, err = getDbConfig(c.cache, c.secretName)
	if err != nil {
		return
	}

	conn, err = c.d.Open(dbConfig.Dsn)
	return
}

func (c *Connector) Driver() driver.Driver {
	return c.d
}

func (a *AWSConfig) InitSqlDB(d driver.Driver) {
	a.d = d
}

func (a *AWSConfig) NewSqlConnector(ctx context.Context) (conn driver.Connector, err error) {
	err = a.isInitialized()
	if err != nil {
		return
	}

	var secretName string
	secretName, err = getSecretName(ctx)
	if err != nil {
		return
	}

	conn = &Connector{
		secretName: secretName,
		d:          a.d,
		cache:      a.cache,
		m:          &sync.Mutex{},
	}
	return
}

func (a *AWSConfig) SqlConnectorProvider() (pName string, p base.SqlConnectorProvider) {
	pName = "aws"
	p = a
	return
}
