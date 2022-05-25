package pgsqldb

import (
	"context"
	"fmt"

	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
	log "github.com/sirupsen/logrus"
)

var (
	pgsqlPoolConnectorProviders map[string]base.PgSqlPoolConnectorProvider
	pgsqlPoolConnectors         map[string]map[string]base.PgSqlPoolConnector
)

func Init(providers ...base.PgSqlPoolConnectorProvider) {
	log.Debug("Initializing local configuration....")
	pgsqlPoolConnectorProviders = make(map[string]base.PgSqlPoolConnectorProvider)
	pgsqlPoolConnectors = make(map[string]map[string]base.PgSqlPoolConnector)

	for _, p := range providers {
		name, provider := p.PgSqlPoolConnectorProvider()
		log.WithFields(log.Fields{
			"ProviderName": name,
		}).Debug("Adding provider")
		pgsqlPoolConnectorProviders[name] = provider
	}
}

func GetPgSqlConnectionPoolConnector(ctx context.Context) (connector base.PgSqlPoolConnector, err error) {
	var (
		cloudVendor string
		provider    base.PgSqlPoolConnectorProvider
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

	if provider, ok = pgsqlPoolConnectorProviders[cloudVendor]; !ok {
		err = fmt.Errorf("there is no cloud provider with the name %s", cloudVendor)
		return
	}

	if connector, ok = pgsqlPoolConnectors[cloudVendor][clientId]; ok {
		log.WithFields(log.Fields{
			"cloudPlatform": cloudVendor,
			"clientId":      clientId,
		}).Debug("postgresql connector retrieved successfully")

		return
	}

	log.WithFields(log.Fields{
		"cloudPlatform": cloudVendor,
		"clientId":      clientId,
	}).Debug("As there is no connector, creating new one")

	connector, err = provider.NewPgSqlPoolConnector(ctx)
	if err != nil {
		return
	}

	pgsqlPoolConnectors[cloudVendor] = make(map[string]base.PgSqlPoolConnector)
	pgsqlPoolConnectors[cloudVendor][clientId] = connector

	log.WithFields(log.Fields{
		"cloudPlatform": cloudVendor,
		"clientId":      clientId,
	}).Debug("postgresql connection pool connector created and retrieved successfully")
	return
}
