package aws

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"sync"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type DbConfig struct {
	Dsn string `mapstructure:"dsn"`
}

type Connector struct {
	secretName string
	cache      *secretcache.Cache
	d          driver.Driver
	m          *sync.Mutex
}

func (c *Connector) Connect(ctx context.Context) (conn driver.Conn, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	result, err := c.cache.GetSecretString(c.secretName)
	if err != nil {
		return
	}

	var dbConfig DbConfig
	err = json.Unmarshal([]byte(result), &dbConfig)
	if err != nil {
		return
	}

	conn, err = c.d.Open(dbConfig.Dsn)
	return
}

func (c *Connector) Driver() driver.Driver {
	return c.d
}
