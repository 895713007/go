package config

import (
	"github.com/mytokenio/go_sdk/config/registry"
	"time"
)

type Config struct {
	Registry registry.Registry
	TTL time.Duration //TODO
}

func NewConfig(opts ...Option) *Config {
	options := newOptions(opts...)
	return &Config{
		Registry: options.Registry,
	}
}

// return string value by key, return error if key not found
// return error if request failed (http registry)
func (c *Config) Get(name string) (string, error) {
	s, err := c.Registry.Get(name)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (c *Config) Set(name string, value string) error {
	return c.Registry.Set(name, value)
}

func (c *Config) String(name string) string {
	s, _ := c.Get(name)
	return s
}

// return string by name, or default value if not found.
func (c *Config) StringOr(name string, dv string) string {
	s, err := c.Registry.Get(name)
	if err != nil {
		return dv
	}

	return s
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
	s, err := c.Registry.Get(name)
	if err != nil {
		return dv
	}

	return toInt(s, dv)
}

func (c *Config) Int64(name string) int64 {
	return c.Int64Or(name, 0)
}

func (c *Config) Int64Or(name string, dv int64) int64 {
	s, err := c.Registry.Get(name)
	if err != nil {
		return dv
	}

	return toInt64(s, dv)
}

func (c *Config) Float64(name string) float64 {
	return c.Float64Or(name, 0)
}

func (c *Config) Float64Or(name string, dv float64) float64 {
	s, err := c.Registry.Get(name)
	if err != nil {
		return dv
	}

	return toFloat64(s, dv)
}

//TODO, more types