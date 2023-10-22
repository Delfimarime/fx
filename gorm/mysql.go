package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlTypedFactory() TypedFactory {
	return TypedFactory{
		Type: MySqlType,
		Factory: func(c config.Config) (*gorm.DB, error) {
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				c.Data.Gorm.Authentication.Username, c.Data.Gorm.Authentication.Password,
				c.Data.Gorm.URL, c.Data.Gorm.Database)
			return gorm.Open(mysql.Open(dsn), &gorm.Config{})
		},
	}
}
