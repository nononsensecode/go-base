package sqldb

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/nononsensecode/go-base"
	"github.com/nononsensecode/go-base/context/ctxtypes"
	"github.com/sirupsen/logrus"
)

var (
	sqlConnProviders map[string]base.SqlConnectorProvider
	sqlConnectors    map[string]map[string]driver.Connector
)

func Init(providers ...base.SqlConnectorProvider) {
	sqlConnProviders = make(map[string]base.SqlConnectorProvider)
	sqlConnectors = make(map[string]map[string]driver.Connector)
	for _, p := range providers {
		name, provider := p.SqlConnectorProvider()
		sqlConnProviders[name] = provider
	}
}

func GetSqlConnector(ctx context.Context) (c driver.Connector, err error) {
	var (
		vendor   string
		provider base.SqlConnectorProvider
		clientId string
		ok       bool
	)
	if vendor, ok = ctx.Value(ctxtypes.CtxVendorKey).(string); !ok {
		err = fmt.Errorf("cloud information is missing")
		return
	}

	if clientId, ok = ctx.Value(ctxtypes.CtxClientIdKey).(string); !ok {
		err = fmt.Errorf("client id information is missing")
		return
	}

	if provider, ok = sqlConnProviders[vendor]; !ok {
		err = fmt.Errorf("there is no cloud provider %s", vendor)
		return
	}

	if c, ok = sqlConnectors[vendor][clientId]; ok {
		logrus.WithFields(logrus.Fields{
			"platform": vendor,
			"clientId": clientId,
		}).Debug("sql connection retrieved successfully")

		return
	}

	logrus.WithFields(logrus.Fields{
		"cloudPlatform": vendor,
		"clientId":      clientId,
	}).Debug("As there is no connection, creating new one")

	c, err = provider.NewSqlConnector(ctx)
	if err != nil {
		return
	}
	sqlConnectors[vendor] = make(map[string]driver.Connector)
	sqlConnectors[vendor][clientId] = c

	logrus.WithFields(logrus.Fields{
		"platform": vendor,
		"clientId": clientId,
	}).Debug("sql connection created and retrieved successfully")

	return
}
