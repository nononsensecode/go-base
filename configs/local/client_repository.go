package local

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	if db == nil {
		panic("client db is nil")
	}

	return &ClientRepository{db}
}

func (c *ClientRepository) GetDsnByClientId(id string) (dsn string, err error) {
	err = c.db.Get(&dsn, "SELECT dsn FROM auth_user_client WHERE id = ?", id)
	return
}

type ClientRepositoryConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
}

func (c *ClientRepositoryConfig) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.DbName)
}
