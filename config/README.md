## Config

### Usage

import package

```
import (
    "github.com/mytokenio/go_sdk/config"
)


c := config.NewConfig()

//return string value by key, return error if key not found
//return error if request failed (http registry)
value, err := c.Get("key")

// shortcut
str := c.String("key")
b := c.Bool("key")
i := c.Int("key")
i64 := c.Int64("key")
f := c.Float64("key")

// default value, if key not exists or request failed
str := c.StringOr("key", "default value")
b := c.BoolOr("key", true)
i := c.IntOr("key", 1234)
i64 := c.Int64Or("key", 43124321)
f := c.Float64Or("key", 2.345)

```

### Registry

support custom config registry, see [http_registry](https://github.com/mytokenio/go_sdk/blob/master/config/registry/http_registry.go) or [mock_registry](https://github.com/mytokenio/go_sdk/blob/master/config/registry/mock_registry.go)

### TODO

more shortcuts support, slice, map, etc.

cache support

