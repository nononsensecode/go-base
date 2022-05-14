package local

import (
	"database/sql"
	"fmt"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	if db == nil {
		panic("client db is nil")
	}

	return &ClientRepository{db}
}

func (c *ClientRepository) GetDsnByClientId(id string) (dsn string, err error) {
	row := c.db.QueryRow("SELECT dsn FROM auth_user_client WHERE id = ?", id)
	err = row.Scan(&dsn)
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
