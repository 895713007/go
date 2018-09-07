## MyToken Go SDK

## Log

log based on `logrus`, support sync to mytoken log server

```
go get github.com/mytokenio/go/log
```

```
import (
    "github.com/mytokenio/go/log"
)

log.Info("xxx log")
log.WithField("kkk", "vvv", "custom_type").Info("log with type & kv data")
```

[more detail](https://github.com/mytokenio/go/tree/master/log)

## Config


```
go get github.com/mytokenio/go/config
```

```
import (
    "github.com/mytokenio/go/config"
)

mc := &MyConfig{}
c := config.NewConfig()

// bind to struct
c.BindTOML(mc)

// or, watch change
c.Watch(func() error {
    err := c.BindTOML(mc)
    // TODO
    return nil
})
```

[more detail](https://github.com/mytokenio/go/tree/master/config)


## TODO

service registry/broker/health

metrics

rpc server/client/protocol

rate limiter, circuit breaker, hytrix

tracing, zipkin or jaeger

cli

queue

demo project

...


