package monitor

import (
	"strings"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
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

func Combine(cfg config.Configuration, monitors ...Monitor) *CombinedMonitor {
	var effectiveMonitors []Monitor
	monitorsByName := make(map[string]Monitor)

	for _, monitor := range monitors {
		monitorsByName[monitor.Name()] = monitor
	}

	for _, name := range cfg.Monitors {
		if monitor, ok := monitorsByName[name]; ok {
			effectiveMonitors = append(effectiveMonitors, monitor)
		}
	}

	enabledMonitorNames := make([]string, 0)
	for _, monitor := range effectiveMonitors {
		enabledMonitorNames = append(enabledMonitorNames, monitor.Name())
	}
	log.Printf("enabled monitors: %s\n", strings.Join(enabledMonitorNames, ","))

	return &CombinedMonitor{effectiveMonitors}
}
