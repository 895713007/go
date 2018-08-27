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
log.WithField("kkk", "vvv").Info("log with kv data")
```

[more detail](https://github.com/mytokenio/go_sdk/tree/master/log)

## Config

...

## TODO

service registry/broker/health

rpc server/client/protocol

...


