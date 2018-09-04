package config

import (
	"testing"
	"runtime"
	"github.com/mytokenio/go_sdk/config/driver"
	"strings"
	"os"
)
const MyConfigJson = `
{
	"api": "http://api.mytokenapi.com",
	"db": {
		"host": "localhost",
		"user": "root",
		"password": "",
		"name": "mytoken"
	},
	"log_servers": ["127.0.0.1:12333", "127.0.0.1:12334"]
}
`
type MyConfig struct {
	API string `json:"api"`
	DB struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"db"`
	LogServers []string `json:"log_servers"`
}

func assert(t *testing.T, actual interface{}, expect interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName = strings.Replace(fileName, wd, "", 1)
	if actual != expect {
		t.Fatalf("expect %v, got %v at (%v:%v)", expect, actual, fileName, line)
	}
}

func newMockConfig() *Config {
	r := driver.NewMockDriver()
	return NewConfig(Driver(r))
}

func TestService(t *testing.T) {
	c := newMockConfig()
	c.Service = "test"
	value := driver.NewValue("mt.service."+c.Service, []byte(MyConfigJson))
	c.Driver.Set(value)

	v, _ := c.GetServiceConfig()
	assert(t, v.String(), MyConfigJson)

	mc := &MyConfig{}
	c.BindJSON(mc)
	assert(t, mc.API, "http://api.mytokenapi.com")
	assert(t, mc.DB.Name, "mytoken")
}

func TestBasic(t *testing.T) {
	c := newMockConfig()
	v, _ := c.Get("foo")
	if v != nil {
		assert(t, "not nil", nil)
	}

	value := driver.NewValue("foo", []byte("bar"))
	c.Driver.Set(value)
	v2, _ := c.Get("foo")
	assert(t, v2.String(), "bar")
}
