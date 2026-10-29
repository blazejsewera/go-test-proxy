package monitor_test

import (
	"errors"
	"testing"

	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/testutils/assert"
)

func TestCombine(t *testing.T) {
	cfg := config.Configuration{Monitors: []string{"counting"}}
	countingMonitor := new(CountingMonitor)
	alwaysOnMonitor := new(AlwaysOnMonitor)

	tested := monitor.Combine(cfg, countingMonitor, alwaysOnMonitor)

	tested.HTTPEvent(event.HTTP{EventType: event.RequestEventType})

	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})
	tested.HTTPEvent(event.HTTP{EventType: event.ResponseEventType})

	tested.Err(errors.New(""))
	tested.Err(errors.New(""))
	tested.Err(errors.New(""))

	assert.Equal(t, 1, countingMonitor.requestsHandled)
	assert.Equal(t, 2, countingMonitor.responsesHandled)
	assert.Equal(t, 3, countingMonitor.errorsHandled)
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

func (m *AlwaysOnMonitor) HTTPEvent(event.HTTP) {
	m.invoked = true
}

func (m *AlwaysOnMonitor) Err(error) {}
