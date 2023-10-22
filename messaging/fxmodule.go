package messaging

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
)

const (
	module = "channel"
)

type ChannelFactory func(channel config.Channel) (Channel, error)

type TypedFactory struct {
	Type    string
	Factory ChannelFactory
}

func From(r ...TypedFactory) func(config.Config) fx.Option {
	return func(c config.Config) fx.Option {
		return New(c, r...)
	}
}

func New(c config.Config, r ...TypedFactory) fx.Option {
	if c.Integrations == nil {
		return fx.Module(module)
	}
	opts := make([]any, 0)
	for key, cfg := range c.Integrations {
		var f ChannelFactory
		for _, info := range r {
			if info.Type == cfg.Type {
				f = info.Factory
				break
			}
		}
		if f == nil {
			f = func(cnf config.Channel) (Channel, error) {
				return nil, fmt.Errorf(`type="%s" not supported for $.integrations["%s"]`, cfg.Type, key)
			}
		}
		opts = append(opts, fx.Annotated{
			Name: key,
			Target: func() (Channel, error) {
				return f(cfg)
			},
		})
	}
	return fx.Module(module, fx.Provide(opts...))
}
