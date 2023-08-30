package proxy_test

import "github.com/blazejsewera/go-test-proxy/proxy"

type MonitorSpy struct {
	Events []proxy.HTTPEvent
}

var _ proxy.Monitor = (*MonitorSpy)(nil)

func (m *MonitorSpy) HTTPEvent(e proxy.HTTPEvent) {
	m.Events = append(m.Events, e)
}

func (m *MonitorSpy) Clear() {
	m.Events = []proxy.HTTPEvent{}
}
