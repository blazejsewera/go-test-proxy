package monitor

import "github.com/blazejsewera/go-test-proxy/event"

type CombinedMonitor struct {
	monitors []Monitor
}

var _ Monitor = (*CombinedMonitor)(nil)

func (c *CombinedMonitor) HTTPEvent(event event.HTTP) {
	for _, m := range c.monitors {
		m.HTTPEvent(event)
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
