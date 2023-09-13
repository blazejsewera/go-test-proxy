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

	var tested proxy.Monitor = monitor.Combine(monitor1, monitor2)
	tested.HTTPEvent(proxy.HTTPEvent{EventType: proxy.RequestEventType})
	tested.HTTPEvent(proxy.HTTPEvent{EventType: proxy.ResponseEventType})
	tested.Err(errors.New(""))

	assert.Equal(t, 1, monitor1.requestsHandled)
	assert.Equal(t, 1, monitor1.responsesHandled)
	assert.Equal(t, 1, monitor1.errorsHandled)
	assert.Equal(t, 1, monitor2.requestsHandled)
	assert.Equal(t, 1, monitor2.responsesHandled)
	assert.Equal(t, 1, monitor2.errorsHandled)
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
