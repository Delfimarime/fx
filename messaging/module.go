package messaging

import (
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
)

func New(c config.Config) fx.Option {
	if c.Channel == nil {
		return fx.Module(module)
	}
	sup := func(name string, info config.Channel) any {
		return fx.Annotated{
			Name:   name,
			Target: GetChannelFrom(info),
		}
	}
	if len(c.Channel) == 1 {
		sup = func(name string, info config.Channel) any {
			return GetChannelFrom(info)
		}
	}
	opts := make([]any, 0)
	for key, info := range c.Channel {
		opts = append(opts, sup(key, info))
	}
	return fx.Module(module, fx.Provide(opts...))
}
