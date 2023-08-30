package proxy_test

import "github.com/blazejsewera/go-test-proxy/proxy"

type MonitorSpy struct {
	Events []proxy.HTTPEvent
	Errors []error
}

var _ proxy.Monitor = (*MonitorSpy)(nil)

func (m *MonitorSpy) HTTPEvent(event proxy.HTTPEvent) {
	m.Events = append(m.Events, event)
}

func (m *MonitorSpy) Err(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *MonitorSpy) Clear() {
	m.Events = []proxy.HTTPEvent{}
	m.Errors = []error{}
}
