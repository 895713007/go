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

## Metrics

```go
//import metrics backend service
import "github.com/mytokenio/go/metrics/logger"

//init with metrics namespace and logger server address
m := logger.New("test", "127.0.0.1:12333")
defer m.Close()

//create counter/gauge instance
c := m.Counter("counter")
g := m.Gauge("test-gauge")

//call counter/gauge
c.Incr(10)
log.Infof("counter value %d", c.Value())

g.Set(1234)
log.Infof("gauge value %d", g.Value())

// with kv pair
c.With("k1": "v1", "k2", "v2", ...).Incr(123)

g.With("k1": "v1", "k2", "v2", ...).Set(123)
```


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


