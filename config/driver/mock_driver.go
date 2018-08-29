package driver

import (
	"sync"
	"github.com/mytokenio/go_sdk/log"
	"fmt"
)

type mockDriver struct {
	sync.RWMutex
	KV  map[string][]byte
}

func NewMockDriver() Driver {
	return &mockDriver{
		KV:  map[string][]byte{},
	}
}

func (c *mockDriver) Get(key string) ([]byte, error) {
	c.RLock()
	v, ok := c.KV[key]
	c.RUnlock()

	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("mock key %s not found", key)
}

func (c *mockDriver) Set(key string, value []byte) error {
	c.Lock()
	c.KV[key] = value
	c.Unlock()
	log.Infof("mock set %s %s", key, value)
	return nil
}

func (c *mockDriver) String() string {
	return "mock"
}