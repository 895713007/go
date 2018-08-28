package registry

import (
	"sync"
	"github.com/mytokenio/go_sdk/log"
	"fmt"
)

type mockRegistry struct {
	sync.RWMutex
	KV  map[string][]byte
}

func NewMockRegistry() Registry {
	return &mockRegistry{
		KV:  map[string][]byte{},
	}
}

func (c *mockRegistry) Get(key string) ([]byte, error) {
	c.RLock()
	v, ok := c.KV[key]
	c.RUnlock()

	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("mock key %s not found", key)
}

func (c *mockRegistry) Set(key string, value []byte) error {
	c.Lock()
	c.KV[key] = value
	c.Unlock()
	log.Infof("mock set %s %s", key, value)
	return nil
}

func (c *mockRegistry) String() string {
	return "mock"
}