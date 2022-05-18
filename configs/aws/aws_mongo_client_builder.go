package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/nononsensecode/go-base"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *AWSConfig) NewMongoDbClientBuilder(ctx context.Context) (builder base.MongoDbClientBuilder, err error) {
	if err = a.isInitialized(); err != nil {
		return
	}

	var secretName string
	if secretName, err = getSecretName(ctx); err != nil {
		return
	}

	builder = &AwsMongoClientBuilder{
		c:          a.cache,
		m:          &sync.Mutex{},
		secretName: secretName,
	}
	return
}

func (a *AWSConfig) MongoDbClientBuilderProvider() (pName string, provider base.MongoDbClientBuilderProvider) {
	pName = configName
	provider = a
	return
}

type AwsMongoClientBuilder struct {
	c          *secretcache.Cache
	m          *sync.Mutex
	secretName string
}

func (b *AwsMongoClientBuilder) GetMongoDbClient(ctx context.Context) (client *mongo.Client, err error) {
	b.m.Lock()
	defer b.m.Unlock()

	var dbConfig DbConfig
	dbConfig, err = getDbConfig(b.c, b.secretName)
	if err != nil {
		return
	}

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.Dsn))
	return
}
