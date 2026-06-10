package proxy

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/interceptor"
)

type Mock struct {
	RoutePattern string
	HandlerFunc  http.HandlerFunc
}

type Builder struct {
	Router            *http.ServeMux
	Monitor           monitor.Monitor
	port              uint16
	mockGroups        map[string][]Mock
	enabledMockGroups []string
}

func NewBuilder() *Builder {
	return &Builder{
		Router:            http.NewServeMux(),
		Monitor:           monitor.NopMonitor{},
		port:              8000,
		mockGroups:        make(map[string][]Mock),
		enabledMockGroups: make([]string, 0),
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

func (b *Builder) WithEnabledMockGroups(groupNames []string) *Builder {
	for _, enabledGroup := range groupNames {
		group, ok := b.mockGroups[enabledGroup]
		if !ok {
			log.Warnf("enable mock group: not found: %s\n", enabledGroup)
			continue
		}
		b.enableMockGroup(enabledGroup, group)
	}

	return b
}

func (b *Builder) enableMockGroup(name string, mocks []Mock) {
	if slices.Contains(b.enabledMockGroups, name) {
		log.Warnf("enable mock group: already enabled: %s\n", name)
		return
	}

	b.enabledMockGroups = append(b.enabledMockGroups, name)

	for _, mock := range mocks {
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

func (b *Builder) WithEnabledAllMocks() *Builder {
	for name, mocks := range b.mockGroups {
		b.enableMockGroup(name, mocks)
	}
	return b
}

func (b *Builder) Build() *http.Server {
	b.logMocks()
	return &http.Server{Addr: b.serverAddr(), Handler: b.Router}
}

func (b *Builder) logMocks() {
	if len(b.enabledMockGroups) == 0 {
		log.Printf("no mocks enabled\n")
		return
	}

	log.Printf("enabled mock groups:\n")
	for _, enabledGroup := range b.enabledMockGroups {
		log.Printf("  %s:\n", enabledGroup)
		mocks := b.mockGroups[enabledGroup]
		for _, mock := range mocks {
			log.Printf("    - %s\n", mock.RoutePattern)
		}
	}
}

func (b *Builder) serverAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", b.port)
}
