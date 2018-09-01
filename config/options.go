package config

import (
	"github.com/mytokenio/go_sdk/config/driver"
	"time"
)

type Option func(*Options)

type Options struct {
	Service   string
	Tags      []string
	TTL       time.Duration //use cache driver for ttl
	Driver    driver.Driver
	OnChanged []func(*Config) error
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Driver: driver.DefaultDriver,
		TTL:    time.Second * 10,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Service(service string) Option {
	return func(o *Options) {
		o.Service = service
	}
}

func Tags(tags []string) Option {
	return func(o *Options) {
		o.Tags = tags
	}
}

func Driver(r driver.Driver) Option {
	return func(o *Options) {
		o.Driver = r
	}
}

func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}

func OnChanged(fn func(c *Config) error) Option {
	return func(o *Options) {
		o.OnChanged = append(o.OnChanged, fn)
	}
}
