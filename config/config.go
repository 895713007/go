package config

import (
	"github.com/mytokenio/go_sdk/config/driver"
	"github.com/mytokenio/go_sdk/log"
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"fmt"
	"time"
)

type Config struct {
	Service string
	Tags    []string
	Driver  driver.Driver
	OnChanged []func(*Config) error
}

func NewConfig(opts ...Option) *Config {
	options := newOptions(opts...)

	//use cache driver if ttl > 0
	if options.TTL > 0 {
		log.Infof("ttl %s, use cache driver", options.TTL.String())
		cacheDriver := driver.NewCacheDriver(
			driver.SubDriver(options.Driver),
			driver.TTL(options.TTL),
		)
		options.Driver = cacheDriver
	}
	c := &Config{
		Service: options.Service,
		Tags: options.Tags,
		Driver:  options.Driver,
		OnChanged: options.OnChanged,
	}
	go c.monitor()
	return c
}

//loop monitor if config changed ?
func (c *Config) monitor() {
	doOnChanged := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("OnChanged callback panic %v", r)
			}
		}()

		for _, fn := range c.OnChanged {
			if err := fn(c); err != nil {
				log.Errorf("OnChanged callback error %s", err)
			}
		}
	}

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
			case <- ticker.C:
				log.Infof("monitoring ...")
				doOnChanged()
		}
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
func (c *Config) BindJSON(obj interface{}) error {
	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, b, err)
		return err
	}
	e := json.Unmarshal(b, obj)
	if e != nil {
		return fmt.Errorf("json unmarshal error %s", e)
	}
	return nil
}

// bind service config to toml struct
func (c *Config) BindTOML(obj interface{}) error {
	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, b, err)
		return err
	}
	e := toml.Unmarshal(b, obj)
	if e != nil {
		return fmt.Errorf("toml unmarshal error %s", e)
	}
	return nil
}
