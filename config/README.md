## Config

### Usage

import package

```
import (
    "github.com/mytokenio/go_sdk/config"
)

```

#### Service

shortcut for:
1. get config raw value by service name
2. bind raw bytes to custom struct via json or toml

custom struct 

```
type MyConfig struct {...}
```


set service name by `config.Service`
```
c := config.NewConfig(config.Service("mt.user"))
```

```
mc := &MyConfig{}
c.BindJSON(mc)
// or
c.BindTOML(mc)
```

code equal to
```
mc := &MyConfig{}
b, _ := c.Get()

json.Unmarshal(b, mc)
// or
toml.Unmarshal(b, mc)
```


#### Get Raw Value by Key

```
//return raw value by key, return error if key not found
//return error if request failed (http registry)
value, err := c.GetKey("key")

// shortcuts for single-value 
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

