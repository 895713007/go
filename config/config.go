package config

import (
	"github.com/mytokenio/go/config/driver"
	"github.com/mytokenio/go/log"
	"time"
	"runtime/debug"
)

const (
	DefaultWatchInterval = 5 * time.Second
	DefaultServicePrefix = "mt.service."
)

type Config struct {
	Service  string
	Tags     []string
	Driver   driver.Driver
	OnChange OnChangeFn
	checkSum string
}

type OnChangeFn func(*Config) error

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
		Tags:    options.Tags,
		Driver:  options.Driver,
	}
	return c
}

func NewFileConfig(path string) *Config {
	return NewConfig(
		Driver(driver.NewFileDriver(driver.Path(path))),
	)
}

func NewHttpConfig(service string, driverOpts ... driver.Option) *Config {
	httpDriver := driver.NewHttpDriver(driverOpts...)

	return NewConfig(
		Service(service),
		Driver(httpDriver),
	)
}

//loop monitor if config changed ?
func (c *Config) Watch(fn OnChangeFn, duration ... time.Duration) {
	c.OnChange = fn
	c.doOnChange()

	interval := DefaultWatchInterval
	if len(duration) > 0 {
		interval = duration[0]
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.watchChange()
			}
		}
	}()
}

func (c *Config) watchChange() {
	v, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("error get config %s %s", c.Service, err)
		return
	}
	log.Infof("watchChange %v %s %s", v, v.CheckSum, c.checkSum)
	if v.CheckSum != c.checkSum {
		err = c.doOnChange()
		if err == nil {
			c.checkSum = v.CheckSum
		}
	}
}

func (c *Config) doOnChange() error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("OnChange callback panic %v %s", r, debug.Stack())
		}
	}()

	if err := c.OnChange(c); err != nil {
		log.Errorf("OnChange callback error %s", err)
		return err
	}
	return nil
}

//======== service config bind to struct =============

// get by service name
func (c *Config) GetServiceConfig() (*driver.Value, error) {
	value, err := c.Get(c.genServiceKey())
	if err == nil && c.checkSum == "" {
		c.checkSum = value.CheckSum
	}
	return value, err
}

//TODO 命名规则目前仅用于 service, 先写死前缀，后续改进
func (c *Config) genServiceKey() string {
	return DefaultServicePrefix + c.Service
}

// bind via config value format
func (c *Config) Bind(obj interface{}) error {
	v, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, v, err)
		return err
	}
	return v.Bind(obj)
}

// bind service config to json struct
func (c *Config) BindJSON(obj interface{}) error {
	v, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, v, err)
		return err
	}
	return v.BindJSON(obj)
}

// bind service config to toml struct
func (c *Config) BindTOML(obj interface{}) error {
	v, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, v, err)
		return err
	}
	return v.BindTOML(obj)
}

// bind service config to yaml struct
func (c *Config) BindYAML(obj interface{}) error {
	v, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf("service config error %s %s %s", c.Service, v, err)
		return err
	}
	return v.BindYAML(obj)
}

// return raw value by key, return error if key not found
// return error if request failed (http driver)
func (c *Config) Get(key string) (*driver.Value, error) {
	v, err := c.Driver.Get(key)
	if err != nil {
		return nil, err
	}
	return v, nil
}
