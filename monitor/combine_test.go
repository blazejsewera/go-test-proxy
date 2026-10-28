package monitor_test

import (
	"errors"
	"testing"

	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/test/assert"
)

func TestCombine(t *testing.T) {
	cfg := config.Configuration{Monitors: []string{"counting"}}
	monitor1 := new(CountingMonitor)
	monitor2 := new(CountingMonitor)
	monitor3 := new(CountingMonitor)
	alwaysOnMonitor := new(AlwaysOnMonitor)

	tested := monitor.Combine(cfg, monitor1, monitor2, alwaysOnMonitor)
	tested.Add(monitor3)

	tested.HTTPEvent(event.HTTP{EventType: event.RequestEventType})

	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})
	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})

	tested.Err(errors.New(""))
	tested.Err(errors.New(""))
	tested.Err(errors.New(""))

	monitors := [...]*CountingMonitor{monitor1, monitor2, monitor3}
	for _, m := range monitors {
		assert.Equal(t, 1, m.requestsHandled)
		assert.Equal(t, 2, m.responsesHandled)
		assert.Equal(t, 3, m.errorsHandled)
	}
	assert.True(t, alwaysOnMonitor.invoked)
}

type CountingMonitor struct {
	requestsHandled  int
	responsesHandled int
	errorsHandled    int
}

func (c *CountingMonitor) Name() string {
	return "counting"
}

func (c *CountingMonitor) HTTPEvent(e event.HTTP) {
	if e.EventType == event.RequestEventType {
		c.requestsHandled++
	} else {
		c.responsesHandled++
	}
}

func (c *CountingMonitor) Err(error) {
	c.errorsHandled++
}

type AlwaysOnMonitor struct {
	invoked bool
}

func (m *AlwaysOnMonitor) Name() string {
	return ""
}

func (m *AlwaysOnMonitor) HTTPEvent(_ event.HTTP) {
	m.invoked = true
}

func (m *AlwaysOnMonitor) Err(error) {}
