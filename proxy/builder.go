package proxy

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
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
	Router                 *http.ServeMux
	Monitor                monitor.Monitor
	targetURL              *url.URL
	port                   uint16
	mockGroups             map[string][]Mock
	enabledMockGroups      []string
	registeredMockPatterns []string
}

func NewBuilder() *Builder {
	return &Builder{
		Router:                 http.NewServeMux(),
		Monitor:                monitor.NopMonitor{},
		port:                   8000,
		mockGroups:             make(map[string][]Mock),
		enabledMockGroups:      make([]string, 0),
		registeredMockPatterns: make([]string, 0),
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

func (b *Builder) WithProxyTarget(rawTargetURL string) *Builder {
	var err error
	b.targetURL, err = url.Parse(rawTargetURL)
	if err != nil {
		log.Fatalf("target host: url parse: %v\n", err)
	}

	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
		reqMonitor.MonitorRequest()

		resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
		internal.ProxyHandler(b.targetURL)(resInterceptor, r)
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
	if _, ok := b.mockGroups[name]; ok {
		log.Warnf("mock group already exists: %s\n", name)
		return b
	}

	b.mockGroups[name] = append(b.mockGroups[name], mocks...)
	return b
}

func (b *Builder) WithEnabledMockGroups(groupNames []string) *Builder {
	for _, groupName := range groupNames {
		group, ok := b.mockGroups[groupName]
		if !ok {
			log.Warnf("enable mock group: not found: %s; skipping\n", groupName)
			continue
		} else if slices.Contains(b.enabledMockGroups, groupName) {
			log.Warnf("enable mock group: already enabled: %s; skipping\n", groupName)
			continue
		}
		b.enabledMockGroups = append(b.enabledMockGroups, groupName)
		b.enableMockGroup(group)
	}

	return b
}

func (b *Builder) enableMockGroup(mocks []Mock) {
	for _, mock := range mocks {
		if slices.Contains(b.registeredMockPatterns, mock.RoutePattern) {
			log.Warnf("enable mock group: route pattern already registered: %s; skipping\n", mock.RoutePattern)
			continue
		}

		wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqMonitor := interceptor.NewRequestMonitor(r, b.Monitor)
			reqMonitor.MonitorRequest()

			resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
			mock.HandlerFunc(resInterceptor, r)
			resInterceptor.MonitorAndForwardResponse()
		})

		b.Router.Handle(mock.RoutePattern, wrapperFunc)
		b.registeredMockPatterns = append(b.registeredMockPatterns, mock.RoutePattern)
	}
}

func (b *Builder) WithEnabledAllMocks() *Builder {
	groupNames := slices.Sorted(maps.Keys(b.mockGroups))
	for _, name := range groupNames {
		mocks := b.mockGroups[name]
		b.enableMockGroup(mocks)
	}
	return b
}

func (b *Builder) Build() *http.ServeMux {
	b.logMocks()
	return b.Router
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
