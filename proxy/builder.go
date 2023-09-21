package proxy

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/interceptor"
	"net/http"
)

type Builder struct {
	Router  *http.ServeMux
	Monitor monitor.Monitor
	port    uint16
}

func NewBuilder() *Builder {
	return &Builder{
		port:    8000,
		Monitor: monitor.NopMonitor{},
		Router:  http.NewServeMux(),
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
	return b.WithHandlerFunc("/", proxyHandler(b.Monitor, url))
}

func (b *Builder) WithHandlerFunc(pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqInterceptor := interceptor.ForRequest(r, b.Monitor)
		reqInterceptor.MonitorRequest()

		resInterceptor := interceptor.NewResponseInterceptor(w, b.Monitor)
		handlerFunc(resInterceptor, r)
		resInterceptor.MonitorAndForwardResponse()
	})

	b.Router.Handle(pattern, wrapperFunc)
	return b
}

func (b *Builder) Build() *http.Server {
	return &http.Server{Addr: b.serverAddr(), Handler: b.Router}
}

func (b *Builder) serverAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", b.port)
}
