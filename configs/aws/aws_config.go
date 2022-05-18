package aws

import (
	"database/sql/driver"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var configName = "aws"

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

func (a *AWSConfig) Init() {
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

func (a *AWSConfig) isInitialized() (err error) {
	if a.cache == nil {
		err = fmt.Errorf("aws configuration is not initialized")
		return
	}
	return
}
