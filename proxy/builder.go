package proxy

import (
	"github.com/blazejsewera/go-test-proxy/urls"
	"io"
	"net/http"
)

type Builder struct {
	Router *http.ServeMux
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithProxyTarget(url string) *Builder {
	if b.Router == nil {
		b.Router = http.NewServeMux()
	}

	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		targetURL := urls.ForwardedURL(url, r.URL)

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

	b.Router.HandleFunc("/", proxyHandler)
	return b
}

func (b *Builder) WithHandlerFunc(pattern string, customHandlerFunc func(w http.ResponseWriter, r *http.Request)) *Builder {
	return b.WithHandler(pattern, http.HandlerFunc(customHandlerFunc))
}

func (b *Builder) WithHandler(pattern string, customHandler http.Handler) *Builder {
	if b.Router == nil {
		b.Router = http.NewServeMux()
	}

	b.Router.Handle(pattern, customHandler)

	return b
}

func (b *Builder) Build() *Server {
	return &Server{
		server: &http.Server{Addr: "", Handler: b.Router},
		router: b.Router,
	}
}
