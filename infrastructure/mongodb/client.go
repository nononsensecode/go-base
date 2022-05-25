package mongodb

import (
	"context"
	"fmt"

	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
	"github.com/sirupsen/logrus"
)

var (
	mongoDbClientBuilderProviders map[string]base.MongoDbClientBuilderProvider
	mongoDbClientBuilders         map[string]map[string]base.MongoDbClientBuilder
)

func Init(providers ...base.MongoDbClientBuilderProvider) {
	mongoDbClientBuilderProviders = make(map[string]base.MongoDbClientBuilderProvider)
	mongoDbClientBuilders = make(map[string]map[string]base.MongoDbClientBuilder)

	for _, p := range providers {
		name, provider := p.MongoDbClientBuilderProvider()
		mongoDbClientBuilderProviders[name] = provider
	}
}

func GetMongoDbClientBuilder(ctx context.Context) (builder base.MongoDbClientBuilder, err error) {
	var (
		cloudVendor string
		provider    base.MongoDbClientBuilderProvider
		clientId    string
		ok          bool
	)

	if cloudVendor, ok = ctx.Value(ctxtypes.CtxVendorKey).(string); !ok {
		err = fmt.Errorf("cloud vendor information is missing")
		return
	}

	if clientId, ok = ctx.Value(ctxtypes.CtxClientIdKey).(string); !ok {
		err = fmt.Errorf("client id is missing")
		return
	}

	if provider, ok = mongoDbClientBuilderProviders[cloudVendor]; !ok {
		err = fmt.Errorf("there is no cloud provider with the name %s", cloudVendor)
		return
	}

	if builder, ok = mongoDbClientBuilders[cloudVendor][clientId]; ok {
		logrus.WithFields(logrus.Fields{
			"cloudPlatform": cloudVendor,
			"clientId":      clientId,
		}).Debug("mongdb client builder already exist and retrieved successfully")

		return
	}

	logrus.WithFields(logrus.Fields{
		"cloudPlatform": cloudVendor,
		"clientId":      clientId,
	}).Debug("As there is no mongodb client builder, creating new one")

	builder, err = provider.NewMongoDbClientBuilder(ctx)
	mongoDbClientBuilders[cloudVendor] = make(map[string]base.MongoDbClientBuilder)
	mongoDbClientBuilders[cloudVendor][clientId] = builder

	logrus.WithFields(logrus.Fields{
		"cloudPlatform": cloudVendor,
		"clientId":      clientId,
	}).Debug("postgresql connection pool created successfully")

	return
}
