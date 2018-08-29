package driver

import (
	"sync"
	"time"
	"github.com/mytokenio/go_sdk/log"
)

type Value struct {
	K         string
	V         []byte
	Timestamp int64
}

type cacheDriver struct {
	sync.RWMutex
	SubDriver Driver
	TTL         time.Duration
	Data        map[string]Value
}

func NewCacheDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	//ttl minimum 5 seconds
	minTTL := time.Second * 5
	if options.TTL > 0 {
		if options.TTL > minTTL {
			minTTL = options.TTL
		} else {
			log.Warnf("minimum ttl 5 seconds")
		}
	}

	if options.SubDriver == nil {
		options.SubDriver = DefaultDriver
	}
	log.Warnf("cache sub driver %s", options.SubDriver.String())
	return &cacheDriver{
		TTL:         minTTL,
		SubDriver: options.SubDriver,
		Data: map[string]Value{},
	}
}

func (c *cacheDriver) Get(key string) ([]byte, error) {
	if cache := c.cacheGet(key); cache != nil {
		return cache, nil
	}

	b, err := c.SubDriver.Get(key)
	if err != nil {
		return nil, err
	}

	c.cacheSet(key, b)

	return b, nil
}

func (c *cacheDriver) Set(key string, value []byte) error {
	c.Lock()
	delete(c.Data, key)
	c.Unlock()

	return c.SubDriver.Set(key, value)
}

func (c *cacheDriver) cacheGet(key string) []byte {
	c.RLock()
	v, ok := c.Data[key]
	c.RUnlock()

	if !ok {
		return nil
	}
	//expired ?
	if v.Timestamp + int64(c.TTL.Seconds()) < time.Now().Unix() {
		c.Lock()
		delete(c.Data, key)
		c.Unlock()

		return nil
	}

	return v.V
}

func (c *cacheDriver) cacheSet(key string, value []byte) error {
	c.Lock()
	defer c.Unlock()

	c.Data[key] = Value{
		K: key,
		V: value,
		Timestamp: time.Now().Unix(),
	}
	return nil
}

func (c *cacheDriver) String() string {
	return "cache"
}