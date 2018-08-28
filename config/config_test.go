package config

import (
	"testing"
	"runtime"
	"github.com/mytokenio/go_sdk/config/registry"
	"strings"
	"os"
)

func assert(t *testing.T, actual interface{}, expect interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName = strings.Replace(fileName, wd, "", 1)
	if actual != expect {
		t.Fatalf("expect %v, got %v at (%v:%v)\n", expect, actual, fileName, line)
	}
}

func newMockConfig() *Config {
	r := registry.NewMockRegistry()
	return NewConfig(Registry(r))
}

func TestString(t *testing.T) {
	c := newMockConfig()
	assert(t, c.String("foo"), "")

	c.Set("foo", "bar")
	assert(t, c.String("foo"), "bar")

	c.Set("foo", "xxx")
	assert(t, c.String("foo"), "xxx")

	assert(t, c.String("not_exists"), "")
	assert(t, c.StringOr("not_exists", "bar"), "bar")
}

func TestBool(t *testing.T) {
	c := newMockConfig()
	c.Set("foo", "true")

	assert(t, c.Bool("foo"), true)
	assert(t, c.BoolOr("foo", false), true)

	c.Set("foo", "False")
	assert(t, c.Bool("foo"), false)
	assert(t, c.BoolOr("foo", true), false)

	assert(t, c.BoolOr("foo_default_value", true), true)
}

func TestInt(t *testing.T) {
	c := newMockConfig()
	c.Set("foo", "123")
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
	c.Set("foo", "123.456")
	assert(t, c.Float64("foo"), float64(123.456))
	assert(t, c.Float64Or("xxx", 234.555), float64(234.555))

	c.Set("foo", "333")
	assert(t, c.Float64("foo"), float64(333))
	assert(t, c.Float64Or("xxx", 234.555), float64(234.555))
}
