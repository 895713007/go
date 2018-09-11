package internal

import (
	"github.com/mytokenio/go/log"
	"github.com/mytokenio/go/metrics"
	"github.com/mytokenio/go/metrics/internal/lv"
)

func NewGauge(name string, obs observeFunc, add observeFunc, labelValues lv.LabelValues) metrics.Gauge {
	return &Gauge{
		name: name,
		lvs:  labelValues,
		obs:  obs,
		add:  add,
	}
}

type Gauge struct {
	name string
	lvs  lv.LabelValues
	obs  observeFunc
	add  observeFunc
	val  valFunc
}

func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		name: g.name,
		lvs:  g.lvs.With(labelValues...),
		obs:  g.obs,
		add:  g.add,
	}
}

func (g *Gauge) Set(value float64) {
	log.Infof("gset %s %f", g.name, value)
	g.obs(g.name, g.lvs, value)
}

func (g *Gauge) Value() float64 {
	return g.val(g.name)
}
