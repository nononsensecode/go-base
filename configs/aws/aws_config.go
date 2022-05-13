package aws

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/nononsensecode/go-base"
)

type AWSConfig struct {
	MaxCacheSize    int   `mapstructure:"maxCacheSize"`
	CacheItemTTL    int64 `mapstructure:"cacheItemTTL"`
	cache           *secretcache.Cache
	d               driver.Driver
	httpMiddlewares []func(http.Handler) http.Handler
}

func (a *AWSConfig) InitDB(d driver.Driver) {
	var err error
	a.d = d
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		panic(err)
	}

	sMgr := secretsmanager.New(sess)
	cacheConfig := secretcache.CacheConfig{
		MaxCacheSize: a.MaxCacheSize,
		CacheItemTTL: a.CacheItemTTL,
	}
	a.cache, err = secretcache.New(
		func(c *secretcache.Cache) { c.CacheConfig = cacheConfig },
		func(c *secretcache.Cache) { c.Client = sMgr },
	)
	if err != nil {
		panic(err)
	}
}

func (a *AWSConfig) NewConnector(ctx context.Context) (conn driver.Connector, err error) {
	if a.cache == nil {
		err = fmt.Errorf("aws configuration is not initialized")
		return
	}

	// secret name has to be created from the client id. Need to write a middleware for
	// guessing secret name and to store it in the context
	var (
		secretName string
		ok         bool
	)
	if secretName, ok = ctx.Value("secretName").(string); !ok {
		err = fmt.Errorf("secret name for aws cannot be found")
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

func (a *AWSConfig) ConnectorProvider() (pName string, p base.ConnectorProvider) {
	pName = "aws"
	p = a
	return
}

// For future use
func (a *AWSConfig) GetMiddlewares() []func(http.Handler) http.Handler {
	return a.httpMiddlewares
}
