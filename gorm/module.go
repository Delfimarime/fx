package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
)

func New(c config.Config) fx.Option {
	if c.Gorm == nil {
		return fx.Module("gorm")
	}
	opts := make([]fx.Option, 0)
	for key, info := range c.Gorm {
		opts = append(opts, fx.Provide(
			NewGormDB(info),
			fx.ParamTags(fmt.Sprintf("name:%s", key)),
		))
	}
	return fx.Module("gorm", opts...)
}
