package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

func getSecretName(ctx context.Context) (s string, err error) {
	// secret name has to be created from the client id. Need to write a middleware for
	// guessing secret name and to store it in the context
	var ok bool

	if s, ok = ctx.Value("secretName").(string); !ok {
		err = fmt.Errorf("secret name for aws cannot be found")
		return
	}

	return
}

func getDbConfig(cache *secretcache.Cache, secretName string) (d DbConfig, err error) {
	result, err := cache.GetSecretString(secretName)
	if err != nil {
		return
	}

	var dbConfig DbConfig
	err = json.Unmarshal([]byte(result), &dbConfig)
	if err != nil {
		return
	}
	return
}
