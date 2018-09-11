package metrics

import "time"

// reference to https://github.com/rcrowley/go-metrics
// internal/kv reference to gokit

var (
	BatchInterval = time.Second * 5
)

// metric service, to create, call, close
type Metrics interface {
	Close() error
	Counter(id string) Counter
	Gauge(id string) Gauge
	String() string
}

type Counter interface {
	Incr(delta int64)
	Decr(delta int64)
	Value() int64
	With(pair ...string) Counter
}

type Gauge interface {
	Set(value int64)
	Value() int64
	With(pair ...string) Gauge
}
