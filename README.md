## MyToken Go SDK

## Log

log based on `logrus`, support sync to mytoken log server

```
go get github.com/mytokenio/go_sdk/log
```

```
import (
    "github.com/mytokenio/go_sdk/log"
)

log.Info("xxx log")
log.WithField("kkk", "vvv", "custom_type").Info("log with type & kv data")
```

[more detail](https://github.com/mytokenio/go_sdk/tree/master/log)

## Config


```
go get github.com/mytokenio/go_sdk/config
```

```
import (
    "github.com/mytokenio/go_sdk/config"
)

mc := &MyConfig{}

c := config.NewConfig(config.Service("mt.user"))
c.BindJSON(mc)
// or
c.BindTOML(mc)

// or custom unmarshal 
```

[more detail](https://github.com/mytokenio/go_sdk/tree/master/config)


## TODO

service registry/broker/health

rpc server/client/protocol

rate limiter, circuit breaker, hytrix

tracing, zipkin or jaeger

cli

queue

demo project

...


