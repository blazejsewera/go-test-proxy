package monitor_test

import (
	"testing"

	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/testutils/assert"
)

func TestCombine(t *testing.T) {
	cfg := config.Configuration{Monitors: []string{"counting", "second"}}
	countingMonitor := new(CountingMonitor)
	secondMonitor := new(SecondMonitor)

	tested := monitor.Combine(cfg, countingMonitor, secondMonitor)

	tested.HTTPEvent(event.HTTP{EventType: event.RequestEventType})

	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})
	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})

	assert.Equal(t, 1, countingMonitor.requestsHandled)
	assert.Equal(t, 2, countingMonitor.responsesHandled)
	assert.True(t, secondMonitor.invoked)
}

type CountingMonitor struct {
	requestsHandled  int
	responsesHandled int
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

type SecondMonitor struct {
	invoked bool
}

func (m *SecondMonitor) Name() string {
	return "second"
}

func (m *SecondMonitor) HTTPEvent(event.HTTP) {
	m.invoked = true
}
