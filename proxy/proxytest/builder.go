package proxytest

import (
	"github.com/blazejsewera/go-test-proxy/urls"
	"io"
	"net/http"
	"net/http/httptest"
)

type TestServerBuilder struct {
	target string
	router *http.ServeMux
}

func Builder() *TestServerBuilder {
	return &TestServerBuilder{}
}

func (b *TestServerBuilder) WithTarget(url string) *TestServerBuilder {
	b.target = url
	return b
}

func (b *TestServerBuilder) WithHandlerFunc(pattern string, customHandlerFunc func(w http.ResponseWriter, r *http.Request)) *TestServerBuilder {
	return b.WithHandler(pattern, http.HandlerFunc(customHandlerFunc))
}

func (b *TestServerBuilder) WithHandler(pattern string, customHandler http.Handler) *TestServerBuilder {
	if b.router == nil {
		b.router = http.NewServeMux()
	}

	b.router.Handle(pattern, customHandler)

	return b
}

func (b *TestServerBuilder) Build() *httptest.Server {
	if b.router == nil {
		b.router = http.NewServeMux()
	}

	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
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
	}

	b.router.HandleFunc("/", proxyHandler)

	return httptest.NewUnstartedServer(b.router)
}
