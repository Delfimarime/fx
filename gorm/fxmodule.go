package gorm

import (
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
)

const module = "gorm"

func New(c config.Config) fx.Option {
	if c.Gorm == nil {
		return fx.Module(module)
	}
	sup := func(name string, info config.GormDatasource) any {
		return fx.Annotated{
			Name:   name,
			Target: GetDatabase(info),
		}
	}
	if len(c.Gorm) == 1 {
		sup = func(name string, info config.GormDatasource) any {
			return GetDatabase(info)
		}
	}
	opts := make([]any, 0)
	for key, info := range c.Gorm {
		opts = append(opts, sup(key, info))
	}
	return fx.Module(module, fx.Provide(opts...))
}
