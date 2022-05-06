package sqldb

import (
	"context"
	"database/sql/driver"
	"fmt"

	"gitlab.com/kaushikayanam/base"
	"gitlab.com/kaushikayanam/base/context/ctxtypes"
)

var (
	connProviders map[string]base.ConnectorProvider
	connectors    map[string]map[string]driver.Connector
)

func Init(providers []base.ConnectorProvider) {
	connProviders = make(map[string]base.ConnectorProvider)
	connectors = make(map[string]map[string]driver.Connector)
	for _, p := range providers {
		name, provider := p.ConnectorProvider()
		connProviders[name] = provider
	}
}

func GetConnector(ctx context.Context) (c driver.Connector, err error) {
	var (
		vendor   string
		provider base.ConnectorProvider
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

	if provider, ok = connProviders[vendor]; !ok {
		err = fmt.Errorf("there is no cloud provider %s", vendor)
		return
	}

	if c, ok = connectors[vendor][clientId]; !ok {
		c, err = provider.NewConnector(ctx)
		if err != nil {
			return
		}
		connectors[vendor] = make(map[string]driver.Connector)
		connectors[vendor][clientId] = c
	}
	return
}