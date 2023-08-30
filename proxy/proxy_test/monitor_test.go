package proxy_test

import "github.com/blazejsewera/go-test-proxy/proxy"

type MonitorSpy struct {
	Events []proxy.HttpEvent
}

var _ proxy.Monitor = MonitorSpy{}

func (m MonitorSpy) HttpEvent(e proxy.HttpEvent) {
	m.Events = append(m.Events, e)
}

func (m MonitorSpy) Clear() {
	m.Events = []proxy.HttpEvent{}
}
