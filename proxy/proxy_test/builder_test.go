package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	*httptest.Server
}

type TestServerBuilder struct {
	builder *proxy.Builder
}

func NewBuilder() *TestServerBuilder {
	return &TestServerBuilder{builder: proxy.NewBuilder()}
}

func (b *TestServerBuilder) WithProxyTarget(url string) *TestServerBuilder {
	b.builder.WithProxyTarget(url)
	return b
}

func (b *TestServerBuilder) WithHandlerFunc(pattern string, customHandlerFunc func(w http.ResponseWriter, r *http.Request)) *TestServerBuilder {
	b.builder.WithHandlerFunc(pattern, customHandlerFunc)
	return b
}

func (b *TestServerBuilder) Build() *TestServer {
	return &TestServer{Server: httptest.NewUnstartedServer(b.builder.Router)}
}
