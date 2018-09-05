package metrics

// reference to https://github.com/rcrowley/go-metrics

type Fields map[string]string

// metric service, to create, call, close
type Metrics interface {
	Init() error
	Close() error
	Counter(id string) Counter
	Gauge(id string) Gauge
	String() string
}

type Counter interface {
	Clear()
	Incr(d uint64)
	Decr(d uint64)
	Value() int64
	WithFields(f Fields) Counter
}

type Gauge interface {
	Clear()
	Update(d int64)
	Value() int64
	WithFields(f Fields) Gauge
}

