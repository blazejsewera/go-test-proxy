package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy"
	"net/http/httptest"
)

func BuildTestServer(builder *proxy.Builder) *httptest.Server {
	return httptest.NewUnstartedServer(builder.Router)
}
