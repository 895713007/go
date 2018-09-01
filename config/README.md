## Config

### Usage

import package

```
import (
    "github.com/mytokenio/go_sdk/config"
)
```

set env for http driver

`CONFIG_SERVER`  `http://xxx.com` 

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
// config with http driver
c := config.NewConfig(
    config.Service("mt.user"),
    config.TTL(time.Second * 10), //cache ttl
    config.Driver(
        driver.NewHttpDriver(
            driver.Host("http://xxx.com"),
            driver.Timeout(time.Second * 3),
        ),
    ),
)
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
b, _ := c.GetServiceConfig()

json.Unmarshal(b, mc)
// or
toml.Unmarshal(b, mc)
```


#### Get Raw Value by Key

```
//return raw value by key, return error if key not found
//return error if request failed (http driver)
value, err := c.Get("key")

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

### Driver

support custom config driver, see [http_driver](https://github.com/mytokenio/go_sdk/blob/master/config/driver/http_driver.go) or [mock_driver](https://github.com/mytokenio/go_sdk/blob/master/config/driver/mock_driver.go)


