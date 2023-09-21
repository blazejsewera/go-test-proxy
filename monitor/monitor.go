package monitor

import "github.com/blazejsewera/go-test-proxy/event"

type Monitor interface {
	HTTPEvent(event event.HTTP)
	Err(err error)
}

type NopMonitor struct{}

var _ Monitor = NopMonitor{}

func (n NopMonitor) HTTPEvent(event.HTTP) {}

func (n NopMonitor) Err(error) {}
