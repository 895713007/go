package registry

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

type cacheRegistry struct {
	sync.RWMutex
	SubRegistry Registry
	TTL         time.Duration
	Data        map[string]Value
}

func NewCacheRegistry(opts ...Option) Registry {
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

	if options.SubRegistry == nil {
		options.SubRegistry = DefaultRegistry
	}
	log.Warnf("cache sub registry %s", options.SubRegistry.String())
	return &cacheRegistry{
		TTL:         minTTL,
		SubRegistry: options.SubRegistry,
		Data: map[string]Value{},
	}
}

func (c *cacheRegistry) Get(key string) ([]byte, error) {
	if cache := c.cacheGet(key); cache != nil {
		return cache, nil
	}

	b, err := c.SubRegistry.Get(key)
	if err != nil {
		return nil, err
	}

	c.cacheSet(key, b)

	return b, nil
}

func (c *cacheRegistry) Set(key string, value []byte) error {
	c.Lock()
	delete(c.Data, key)
	c.Unlock()

	return c.SubRegistry.Set(key, value)
}

func (c *cacheRegistry) cacheGet(key string) []byte {
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

func (c *cacheRegistry) cacheSet(key string, value []byte) error {
	c.Lock()
	defer c.Unlock()

	c.Data[key] = Value{
		K: key,
		V: value,
		Timestamp: time.Now().Unix(),
	}
	return nil
}

func (c *cacheRegistry) String() string {
	return "cache"
}