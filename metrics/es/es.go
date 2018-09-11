package es

import (
	"github.com/mytokenio/go/log"
	"github.com/mytokenio/go/metrics"
	"github.com/mytokenio/go/metrics/internal"
	"github.com/mytokenio/go/metrics/internal/lv"
	"io"
	"net"
	"time"
)

type es struct {
	namespace string
	addr      string
	exit      chan bool
	counters  *lv.Space
	gauges    *lv.Space
}

func New(namespace, addr string) metrics.Metrics {
	m := &es{
		namespace: namespace,
		addr:      addr,
		exit:      make(chan bool),
		counters:  lv.NewSpace(),
		gauges:    lv.NewSpace(),
	}
	go m.run()
	return m
}

func (m *es) String() string {
	return "es"
}

func (m *es) Counter(name string) metrics.Counter {
	return internal.NewCounter(name, m.counters.Observe, lv.LabelValues{})
}

func (m *es) counterVal(name string) float64 {
	var val float64
	m.counters.WalkNode(name, func(lvs lv.LabelValues, values []float64) bool {
		val += sum(values)
		return true
	})
	return val
}

func (m *es) Gauge(name string) metrics.Gauge {
	return internal.NewGauge(name, m.gauges.Observe, m.gauges.Add, lv.LabelValues{})
}

func (m *es) gaugeVal(name string) float64 {
	var val float64
	m.counters.WalkNode(name, func(lvs lv.LabelValues, values []float64) bool {
		val = last(values)
		return true
	})
	return val
}

func (m *es) run() {
	t := time.NewTicker(metrics.BatchInterval)
	conn, err := net.DialTimeout("udp", m.addr, time.Second)
	if err != nil {
		log.Errorf("failed dial log server %v", m.addr)
	}
	defer conn.Close()

	for {
		select {
		case <-m.exit:
			t.Stop()
			return
		case <-t.C:
			m.writeTo(conn)
		}
	}
}

func (m *es) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}

func (m *es) writeTo(w io.Writer) (count int64, err error) {
	var n int

	m.counters.Reset().Walk(func(name string, lvs lv.LabelValues, values []float64) bool {
		n, err = w.Write(newCounterMessage(name, sum(values), lvPair(lvs)))
		if err != nil {
			return false
		}

		count += int64(n)
		return true
	})
	if err != nil {
		return count, err
	}

	m.gauges.Reset().Walk(func(name string, lvs lv.LabelValues, values []float64) bool {
		n, err = w.Write(newGaugeMessage(name, last(values), lvPair(lvs)))
		if err != nil {
			return false
		}

		count += int64(n)
		return true
	})
	if err != nil {
		return count, err
	}

	return count, err
}

func lvPair(labelValues []string) map[string]string {
	if len(labelValues)%2 != 0 {
		panic("lvPair received a labelValues with an odd number of strings")
	}
	ret := make(map[string]string, len(labelValues)/2)
	for i := 0; i < len(labelValues); i += 2 {
		ret[labelValues[i]] = labelValues[i+1]
	}
	return ret
}

func sum(a []float64) float64 {
	var v float64
	for _, f := range a {
		v += f
	}
	return v
}

func last(a []float64) float64 {
	return a[len(a)-1]
}
