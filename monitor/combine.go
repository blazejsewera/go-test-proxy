package monitor

import (
	"slices"

	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/event"
)

type CombinedMonitor struct {
	monitors []Monitor
}

var _ Monitor = (*CombinedMonitor)(nil)

func (c *CombinedMonitor) Name() string {
	return "combined"
}

func (c *CombinedMonitor) HTTPEvent(e event.HTTP) {
	for _, m := range c.monitors {
		m.HTTPEvent(e)
	}
}

func (c *CombinedMonitor) Err(err error) {
	for _, m := range c.monitors {
		m.Err(err)
	}
}

func (c *CombinedMonitor) Add(m Monitor) {
	c.monitors = append(c.monitors, m)
}

func Combine(cfg config.Configuration, monitors ...Monitor) *CombinedMonitor {
	var effectiveMonitors []Monitor
	for _, m := range monitors {
		if slices.Contains(cfg.Monitors, m.Name()) || m.Name() == "" {
			effectiveMonitors = append(effectiveMonitors, m)
		}
	}

	return &CombinedMonitor{effectiveMonitors}
}
