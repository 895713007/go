package driver

import (
	"sync"
	"github.com/mytokenio/go/log"
	"fmt"
	"errors"
)

type mockDriver struct {
	sync.RWMutex
	KV  map[string]*Value
}

func NewMockDriver() Driver {
	return &mockDriver{
		KV:  map[string]*Value{},
	}
}

func (d *mockDriver) List() ([]*Value, error) {
	return nil, errors.New("TODO")
}

func (d *mockDriver) Get(key string) (*Value, error) {
	d.RLock()
	v, ok := d.KV[key]
	d.RUnlock()

	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("mock key %s not found", key)
}

func (d *mockDriver) Set(value *Value) error {
	d.Lock()
	d.KV[value.K] = value
	d.Unlock()
	log.Infof("mock set %s %s", value.K, value)
	return nil
}

func (d *mockDriver) String() string {
	return "mock"
}