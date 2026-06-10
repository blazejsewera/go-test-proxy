package proxy

import (
	"fmt"
	"net/http"

	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/interceptor"
)

type Mock struct {
	RoutePattern string
	HandlerFunc  http.HandlerFunc
}

type Builder struct {
	Router     *http.ServeMux
	Monitor    monitor.Monitor
	port       uint16
	mockGroups map[string][]Mock
}

func NewBuilder() *Builder {
	return &Builder{
		Router:     http.NewServeMux(),
		Monitor:    monitor.NopMonitor{},
		port:       8000,
		mockGroups: make(map[string][]Mock),
	}
}

func (b *Builder) WithPort(port uint16) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithMonitor(monitor monitor.Monitor) *Builder {
	b.Monitor = monitor
	return b
}

func (b *Builder) WithProxyTarget(url string) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
		reqMonitor.MonitorRequest()

		resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
		internal.ProxyHandler(b.Monitor, url)(resInterceptor, r)
		resInterceptor.MonitorAndForwardResponse()
	})

	b.Router.Handle("/", wrapperFunc)
	return b
}

func (b *Builder) WithMock(pattern string, handlerFunc http.HandlerFunc) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
		reqMonitor.MonitorRequest()

		resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
		handlerFunc(resInterceptor, r)
		resInterceptor.MonitorAndForwardResponse()
	})

	b.Router.Handle(pattern, wrapperFunc)
	return b
}

func (b *Builder) WithMockGroup(name string, mocks ...Mock) *Builder {
	b.mockGroups[name] = append(b.mockGroups[name], mocks...)
	return b
}

func (b *Builder) WithEnabledMockGroups(groupNames ...string) *Builder {
	for _, enabledGroup := range groupNames {
		group, ok := b.mockGroups[enabledGroup]
		if !ok {
			continue
		}
		for _, mock := range group {
			wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
				reqMonitor.MonitorRequest()

				resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
				mock.HandlerFunc(resInterceptor, r)
				resInterceptor.MonitorAndForwardResponse()
			})

			b.Router.Handle(mock.RoutePattern, wrapperFunc)
		}
	}

	return b
}

func (b *Builder) WithEnabledAllMocks() *Builder {
	for _, mockGroup := range b.mockGroups {
		for _, mock := range mockGroup {
			wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
				reqMonitor.MonitorRequest()

				resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
				mock.HandlerFunc(resInterceptor, r)
				resInterceptor.MonitorAndForwardResponse()
			})

			b.Router.Handle(mock.RoutePattern, wrapperFunc)
		}
	}
	return b
}

func (b *Builder) Build() *http.Server {
	return &http.Server{Addr: b.serverAddr(), Handler: b.Router}
}

func (b *Builder) serverAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", b.port)
}
