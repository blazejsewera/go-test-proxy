package monitor_test

import (
	"errors"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"testing"
)

func TestCombine(t *testing.T) {
	monitor1 := new(CountingMonitor)
	monitor2 := new(CountingMonitor)
	monitor3 := new(CountingMonitor)

	tested := monitor.Combine(monitor1, monitor2)
	tested.Add(monitor3)

	tested.HTTPEvent(proxy.HTTPEvent{EventType: proxy.RequestEventType})

	tested.HTTPEvent(proxy.HTTPEvent{EventType: proxy.ResponseEventType})
	tested.HTTPEvent(proxy.HTTPEvent{EventType: proxy.ResponseEventType})

	tested.Err(errors.New(""))
	tested.Err(errors.New(""))
	tested.Err(errors.New(""))

	monitors := [...]*CountingMonitor{monitor1, monitor2, monitor3}
	for _, m := range monitors {
		assert.Equal(t, 1, m.requestsHandled)
		assert.Equal(t, 2, m.responsesHandled)
		assert.Equal(t, 3, m.errorsHandled)
	}
}

type CountingMonitor struct {
	requestsHandled  int
	responsesHandled int
	errorsHandled    int
}

func (c *CountingMonitor) HTTPEvent(event proxy.HTTPEvent) {
	if event.EventType == proxy.RequestEventType {
		c.requestsHandled++
	} else {
		c.responsesHandled++
	}
}

func (c *CountingMonitor) Err(error) {
	c.errorsHandled++
}
