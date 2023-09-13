package monitor

import "github.com/blazejsewera/go-test-proxy/proxy"

type CombinedMonitor struct {
	monitors []proxy.Monitor
}

var _ proxy.Monitor = (*CombinedMonitor)(nil)

func (c *CombinedMonitor) HTTPEvent(event proxy.HTTPEvent) {
	for _, m := range c.monitors {
		m.HTTPEvent(event)
	}
}

func (c *CombinedMonitor) Err(err error) {
	for _, m := range c.monitors {
		m.Err(err)
	}
}

func (c *CombinedMonitor) Add(m proxy.Monitor) {
	c.monitors = append(c.monitors, m)
}

func Combine(monitors ...proxy.Monitor) *CombinedMonitor {
	return &CombinedMonitor{monitors}
}
