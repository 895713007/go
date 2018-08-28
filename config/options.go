package config

import (
	"github.com/mytokenio/go_sdk/config/registry"
	"time"
)

type Option func(*Options)

type Options struct {
	Registry registry.Registry
	TTL time.Duration
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Registry:  registry.DefaultRegistry,
		TTL: time.Second * 10,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}