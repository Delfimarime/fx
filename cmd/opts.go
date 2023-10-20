package cmd

import (
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
)

func NewOpts(seq ...func(*Opts)) Opts {
	opts := Opts{
		api:       "N/A",
		options:   make([]fx.Option, 0),
		factories: make([]func(config config.Config) fx.Option, 0),
	}
	for _, e := range seq {
		e(&opts)
	}
	return opts
}

type Opts struct {
	api       string
	options   []fx.Option
	factories []func(config config.Config) fx.Option
}

func WithApi(name string) func(*Opts) {
	return func(opts *Opts) {
		opts.api = name
	}
}

func WithOption(option fx.Option) func(*Opts) {
	return func(opts *Opts) {
		if opts.options == nil {
			opts.options = make([]fx.Option, 0)
		}
		opts.options = append(opts.options, option)
	}
}

func WithOptionFactory(f func(config config.Config) fx.Option) func(*Opts) {
	return func(opts *Opts) {
		if opts.factories == nil {
			opts.factories = make([]func(config config.Config) fx.Option, 0)
		}
		opts.factories = append(opts.factories, f)
	}
}
