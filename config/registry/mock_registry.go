package registry

import (
	"time"
	"errors"
	"sync"
)

type mockRegistry struct {
	TTL time.Duration
	KV  map[string]string
	sync.RWMutex
}

func NewMockRegistry(opts ...Option) Registry {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	//ttl minimum 5 seconds
	minTTL := time.Second * 5
	if options.TTL > minTTL {
		minTTL = options.TTL
	}

	return &mockRegistry{
		TTL: minTTL,
		KV:  map[string]string{},
	}
}

func (c *mockRegistry) Get(name string) (string, error) {
	c.RLock()
	v, ok := c.KV[name]
	c.RUnlock()

	if ok {
		return v, nil
	}
	return "", errors.New("not found")
}

func (c *mockRegistry) Set(name string, value string) error {
	c.Lock()
	c.KV[name] = value
	c.Unlock()
	return nil
}
