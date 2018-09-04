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

#### Service Config

custom struct 

```
type MyConfig struct {...}
```


set service name by `config.Service`
```
// config with http driver
c := config.NewConfig(
    config.Service("user"),
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


### Watch Change

default watch interval 5 seconds,
```
c.Watch(func(c *config.Config) error {
    err := c.BindTOML(mc)
    if err != nil {
        log.Errorf("config bind error %s", err)
        return err
    }

    log.Infof("service config changed %v", mc)
    return nil
})
```

you can pass second parameter to control interval

`c.Watch(callback, 10 * time.Second)`


### File Driver

the default config driver, default file name `config.toml`
```
c := config.NewConfig()
c.BindTOML(...)
// or
c.Watch(...)
```

### UI (for http driver)

```
cd ui/
go build .
./ui
```
open `http://127.0.0.1:5556`

### Other

TODO


