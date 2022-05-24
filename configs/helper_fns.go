package configs

import "github.com/nononsensecode/go-base/infrastructure/sqldb"

func (c Config) IsSqlEnabled() bool {
	return c.Server.Persistence.SqlEnable
}

func (c Config) IsPgxEnabled() bool {
	return c.Server.Persistence.PgxEnable
}

func (c Config) IsMongoDbEnabled() bool {
	return c.Server.Persistence.MongoEnable
}

func (c Config) SqlDbType() sqldb.DbType {
	return c.Server.Persistence.Sql.SqlDbType()
}

func (c Config) HttpAddress() string {
	return c.Server.Http.Address()
}

func (c Config) HttpApiPrefix() string {
	return c.Server.Http.ApiPrefix
}

func (c Config) HttpCorsOrigins() []string {
	return c.Server.Http.CorsAllowedOrigins
}
