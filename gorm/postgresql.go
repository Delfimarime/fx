package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresqlTypedFactory() TypedFactory {
	return TypedFactory{
		Type: PostgresSqlType,
		Factory: func(c config.Config) (*gorm.DB, error) {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				c.Data.Gorm.Host, c.Data.Gorm.Authentication.Username,
				c.Data.Gorm.Authentication.Password, c.Data.Gorm.Database)
			return gorm.Open(postgres.Open(dsn), &gorm.Config{})
		},
	}
}
