package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	MySqlType       = "mysql"
	SqlLiteType     = "sqlite"
	MariadbType     = "mariadb"
	PostgresSqlType = "postgres"
)

func NewGormDB(info config.GormDatasource) func() (*gorm.DB, error) {
	return func() (*gorm.DB, error) {
		var d gorm.Dialector
		switch info.Type {
		case MySqlType, MariadbType:
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", info.Authentication.Username, info.Authentication.Password, info.URL, info.Database)
			d = mysql.Open(dsn)
		case PostgresSqlType:
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", info.Host, info.Authentication.Username, info.Authentication.Password, info.Database)
			d = postgres.Open(dsn)
		case SqlLiteType:
			d = sqlite.Open(info.URL)
		default:
			return nil, fmt.Errorf("unsupported database type: %s", info.Type)
		}
		return gorm.Open(d, &gorm.Config{})
	}
}
