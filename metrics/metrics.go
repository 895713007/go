package metrics

import "time"

// reference to https://github.com/rcrowley/go-metrics
// internal/kv reference to gokit

var (
	BatchInterval = time.Second * 2
)

// metric service, to create, call, close
type Metrics interface {
	Close() error
	Counter(id string) Counter
	Gauge(id string) Gauge
	String() string
}

type Counter interface {
	Incr(d float64)
	Decr(d float64)
	Value() float64
	With(pair ...string) Counter
}

type Gauge interface {
	Set(d float64)
	Value() float64
	With(pair ...string) Gauge
}
