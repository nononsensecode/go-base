package aws

import (
	"database/sql/driver"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/nononsensecode/go-base/infrastructure/sqldb"
	"github.com/nononsensecode/go-base/interfaces/httpsrvr"
)

const (
	ConfigName = "aws"
)

type AWSConfig struct {
	MaxCacheSize    int   `mapstructure:"maxCacheSize"`
	CacheItemTTL    int64 `mapstructure:"cacheItemTTL"`
	cache           *secretcache.Cache
	d               driver.Driver
	httpMiddlewares []func(http.Handler) http.Handler
}

type DbConfig struct {
	Dsn string `mapstructure:"dsn"`
}

func (a *AWSConfig) Init(sqlDriver driver.Driver) (err error) {
	if sqlDriver == nil {
		err = fmt.Errorf("sql driver is nil")
		return
	}

	if a.cache == nil {
		err = fmt.Errorf("aws configuration is not initialized")
		return
	}

	var sess *session.Session
	sess, err = session.NewSession(&aws.Config{})
	if err != nil {
		return
	}

	a.d = sqlDriver

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
		return
	}

	sqldb.Init(a)
	httpsrvr.AddMiddlewares(a.GetMiddlewares()...)

	return
}
