package internal

import (
	"github.com/mytokenio/go/metrics"
	"github.com/mytokenio/go/metrics/internal/lv"
)

func NewCounter(name string, obs observeFunc, labelValues lv.LabelValues) metrics.Counter {
	return &Counter{
		name: name,
		lvs:  labelValues,
		obs:  obs,
	}
}

type Counter struct {
	name string
	lvs  lv.LabelValues
	obs  observeFunc
	val  valFunc
}

func (c *Counter) With(labelValues ...string) metrics.Counter {
	return &Counter{
		name: c.name,
		lvs:  c.lvs.With(labelValues...),
		obs:  c.obs,
	}
}

func (c *Counter) Incr(delta int64) {
	c.obs(c.name, c.lvs, float64(delta))
}

func (c *Counter) Decr(delta int64) {
	c.obs(c.name, c.lvs, float64(-delta))
}

func (c *Counter) Value() int64 {
	return int64(c.val(c.name))
}
