package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const (
	module = "gorm"
)

const (
	MySqlType       = "mysql"
	SqlLiteType     = "sqlite"
	MariadbType     = "mariadb"
	PostgresSqlType = "postgres"
)

type DbFactory func(config config.Config) (*gorm.DB, error)

type TypedFactory struct {
	Type    string
	Factory DbFactory
}

func From(r ...TypedFactory) func(config.Config) fx.Option {
	return func(c config.Config) fx.Option {
		return New(c, r...)
	}
}

func New(c config.Config, r ...TypedFactory) fx.Option {
	if c.Data.Gorm == nil {
		return fx.Module(module)
	}
	var f DbFactory
	for _, info := range r {
		if info.Type == c.Data.Gorm.Type {
			f = info.Factory
			break
		}
	}
	if f == nil {
		f = func(config config.Config) (*gorm.DB, error) {
			return nil, fmt.Errorf(`$.data.gorm.type="%s" not supported`, c.Data.Gorm.Type)
		}
	}
	return fx.Module(module, fx.Provide(f))
}
