package proxytest

import (
	"github.com/blazejsewera/go-test-proxy/urls"
	"io"
	"net/http"
	"net/http/httptest"
)

type TestServerBuilder struct {
	target string
	router http.ServeMux
}

func (b *TestServerBuilder) WithTarget(url string) *TestServerBuilder {
	b.target = url
	return b
}

func (b *TestServerBuilder) Build() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL := urls.ForwardedURL(b.target, r.URL)

		r.Host = targetURL.Host
		r.RequestURI = ""
		r.URL = targetURL
		r.Header.Clone()

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			return
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			return
		}
	})

	return httptest.NewUnstartedServer(handler)
}

func (b *TestServerBuilder) WithHandler(path string, customHandler http.HandlerFunc) *TestServerBuilder {
	panic("not yet implemented")
}

func Builder() *TestServerBuilder {
	return &TestServerBuilder{}
}
