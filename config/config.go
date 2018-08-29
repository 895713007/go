package config

import (
	"github.com/mytokenio/go_sdk/config/registry"
	"github.com/mytokenio/go_sdk/log"
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Service  string
	Registry registry.Registry
}

func NewConfig(opts ...Option) *Config {
	options := newOptions(opts...)

	//use cache registry if ttl > 0
	if options.TTL > 0 {
		log.Infof("ttl %s, use cache registry", options.TTL.String())
		cacheRegistry := registry.NewCacheRegistry(
			registry.SubRegistry(options.Registry),
			registry.TTL(options.TTL),
		)
		options.Registry = cacheRegistry
	}
	return &Config{
		Service:  options.Service,
		Registry: options.Registry,
	}
}

//======== service config bind to struct =============

// get by service name
func (c *Config) GetServiceConfig() ([]byte, error) {
	if c.Service == "" {
		return nil, errors.New("service name not set")
	}
	return c.Get(c.Service)
}

// bind service config to json struct
func (c *Config) BindJSON(obj interface{}) {
	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s", c.Service, b)
		return
	}
	e := json.Unmarshal(b, obj)
	if e != nil {
		log.Errorf("json unmarshal error %s", e)
	}
}

// bind service config to toml struct
func (c *Config) BindTOML(obj interface{}) {
	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s", c.Service, b)
		return
	}
	e := toml.Unmarshal(b, obj)
	if e != nil {
		log.Errorf("toml unmarshal error %s", e)
	}
}

//======== shortcuts for single-value key ==============

// return raw value by key, return error if key not found
// return error if request failed (http registry)
func (c *Config) Get(key string) ([]byte, error) {
	b, err := c.Registry.Get(key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// return string
func (c *Config) String(key string) string {
	return c.StringOr(key, "")
}

// return string by name, or default value if not found.
func (c *Config) StringOr(key string, dv string) string {
	b, err := c.Get(key)
	if err != nil {
		return dv
	}

	return toString(b, "")
}

func (c *Config) Bool(name string) bool {
	return c.BoolOr(name, false)
}

func (c *Config) BoolOr(name string, dv bool) bool {
	s, err := c.Registry.Get(name)
	if err != nil {
		return dv
	}

	return toBool(s, dv)
}

func (c *Config) Int(name string) int {
	return c.IntOr(name, 0)
}

func (c *Config) IntOr(name string, dv int) int {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toInt(b, dv)
}

func (c *Config) Int64(name string) int64 {
	return c.Int64Or(name, 0)
}

func (c *Config) Int64Or(name string, dv int64) int64 {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toInt64(b, dv)
}

func (c *Config) Float64(name string) float64 {
	return c.Float64Or(name, 0)
}

func (c *Config) Float64Or(name string, dv float64) float64 {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toFloat64(b, dv)
}

//TODO, more types
