package gorm

import (
	"github.com/delfimarime/fx/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteTypedFactory() TypedFactory {
	return TypedFactory{
		Type: SqlLiteType,
		Factory: func(c config.Config) (*gorm.DB, error) {
			return gorm.Open(sqlite.Open(c.Data.Gorm.URL), &gorm.Config{})
		},
	}
}
