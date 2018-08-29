package config

import (
	"testing"
	"runtime"
	"github.com/mytokenio/go_sdk/config/driver"
	"strings"
	"os"
	"github.com/mytokenio/go_sdk/log"
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
		t.Fatalf("expect %v, got %v at (%v:%v)\n", expect, actual, fileName, line)
	}
}

func newMockConfig() *Config {
	r := driver.NewMockDriver()
	return NewConfig(Driver(r))
}

func TestService(t *testing.T) {
	c := newMockConfig()
	c.Service = "test.service.name"
	c.Driver.Set(c.Service, []byte(MyConfigJson))

	b, _ := c.GetServiceConfig()
	assert(t, string(b), MyConfigJson)

	mc := &MyConfig{}
	c.BindJSON(mc)
	assert(t, mc.API, "http://api.mytokenapi.com")
	assert(t, mc.DB.Name, "mytoken")
}

func TestBasic(t *testing.T) {
	c := newMockConfig()
	b, _ := c.Get("foo")
	assert(t, string(b), "")

	c.Driver.Set("foo", []byte("bar"))
	b2, _ := c.Get("foo")
	assert(t, string(b2), "bar")
}

func TestString(t *testing.T) {
	c := newMockConfig()
	assert(t, c.String("foo"), "")

	c.Driver.Set("foo", []byte("bar"))
	assert(t, c.String("foo"), "bar")

	c.Driver.Set("foo", []byte("xxx"))
	log.Infof("driver name %s", c.Driver.String())
	assert(t, c.String("foo"), "xxx")

	assert(t, c.String("not_exists"), "")
	assert(t, c.StringOr("not_exists", "bar"), "bar")
}

func TestBool(t *testing.T) {
	c := newMockConfig()
	c.Driver.Set("foo", []byte("true"))

	assert(t, c.Bool("foo"), true)
	assert(t, c.BoolOr("foo", false), true)

	c.Driver.Set("foo", []byte("False"))
	assert(t, c.Bool("foo"), false)
	assert(t, c.BoolOr("foo", true), false)

	assert(t, c.BoolOr("foo_default_value", true), true)
}

func TestInt(t *testing.T) {
	c := newMockConfig()
	c.Driver.Set("foo", []byte("123"))
	assert(t, c.Int("foo"), 123)
	assert(t, c.IntOr("foo", 222), 123)

	assert(t, c.Int("xx"), 0)
	assert(t, c.IntOr("xx", 123), 123)

	assert(t, c.Int64Or("foo", 222), int64(123))
	assert(t, c.Int64Or("xx", 222), int64(222))
	assert(t, c.Int64("xx"), int64(0))
}

func TestFloat(t *testing.T) {
	c := newMockConfig()
	c.Driver.Set("foo", []byte("123.456"))
	assert(t, c.Float64("foo"), float64(123.456))
	assert(t, c.Float64Or("xxx", 234.555), float64(234.555))

	c.Driver.Set("foo", []byte("333"))
	assert(t, c.Float64("foo"), float64(333))
	assert(t, c.Float64Or("xxx", 234.555), float64(234.555))
}
