package registry

import (
	"time"
)

type Registry interface {
	Get(string) (string, error)
	Set(string, string) error
}

var (
	DefaultRegistry = NewHttpRegistry()
)

type Option func(*Options)

type Options struct {
	Host string
	TTL time.Duration
	Timeout time.Duration
}

// Host is the registry addresse to use
func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}
