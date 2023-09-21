package monitor

import (
	"github.com/blazejsewera/go-test-proxy/monitor/event"
)

type CombinedMonitor struct {
	monitors []Monitor
}

var _ Monitor = (*CombinedMonitor)(nil)

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

func Combine(monitors ...Monitor) *CombinedMonitor {
	return &CombinedMonitor{monitors}
}
