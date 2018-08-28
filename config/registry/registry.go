package registry

import (
	"time"
)

type Registry interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	String() string
}

var (
	DefaultRegistry = NewHttpRegistry()
)

type Option func(*Options)

type Options struct {
	Host string     //for http registry
	Timeout time.Duration
	SubRegistry Registry //for cache registry
	TTL time.Duration
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

func SubRegistry(reg Registry) Option {
	return func(o *Options) {
		o.SubRegistry = reg
	}
}
