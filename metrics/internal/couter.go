package internal

import (
	"github.com/mytokenio/go/log"
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

func (c *Counter) Incr(delta float64) {
	log.Infof("incr %s %f", c.name, delta)
	c.obs(c.name, c.lvs, delta)
}

func (c *Counter) Decr(delta float64) {
	c.obs(c.name, c.lvs, -delta)
}

func (c *Counter) Value() float64 {
	return c.val(c.name)
}
