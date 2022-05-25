package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LocalMongoDbClientBuilder struct {
	r        ClientRepository
	m        *sync.Mutex
	clientId string
}

func (mb *LocalMongoDbClientBuilder) GetMongoDbClient(ctx context.Context) (client *mongo.Client, err error) {
	mb.m.Lock()
	defer mb.m.Unlock()

	var dsn string
	dsn, err = mb.r.GetDsnByClientId(mb.clientId)
	if err != nil {
		return
	}

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	return
}

func (l *LocalConfig) NewMongoDbClientBuilder(ctx context.Context) (builder base.MongoDbClientBuilder, err error) {
	var (
		clientId string
		ok       bool
	)

	if clientId, ok = ctx.Value(ctxtypes.CtxClientIdKey).(string); !ok {
		err = fmt.Errorf("client id is missing in the context")
		return
	}

	builder = &LocalMongoDbClientBuilder{
		r:        *l.clientRepo,
		m:        &sync.Mutex{},
		clientId: clientId,
	}
	return
}

func (l *LocalConfig) MongoDbClientBuilderProvider() (pName string, provider base.MongoDbClientBuilderProvider) {
	pName = ConfigName
	provider = l
	return
}
