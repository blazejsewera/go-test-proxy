package proxy

import (
	"bytes"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/urls"
	"io"
	"net/http"
)

type Builder struct {
	Router  *http.ServeMux
	Monitor Monitor
	port    uint16
}

func NewBuilder() *Builder {
	return &Builder{
		port:    8080,
		Monitor: NoopMonitor{},
		Router:  http.NewServeMux(),
	}
}

func (b *Builder) WithPort(port uint16) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithMonitor(monitor Monitor) *Builder {
	b.Monitor = monitor
	return b
}

func (b *Builder) WithProxyTarget(url string) *Builder {
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

	return b.WithHandlerFunc("/", proxyHandler)
}

func (b *Builder) WithHandlerFunc(pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.Monitor.HTTPEvent(requestHTTPEvent(r))
		handlerFunc(w, r)
	})
	b.Router.Handle(pattern, wrapperFunc)
	return b
}

func requestHTTPEvent(r *http.Request) HTTPEvent {
	body, bodyReader := bodyToStringAndReader(r.Body)
	r.Body = bodyReader
	return HTTPEvent{
		EventType: RequestEventType,
		Header:    copyHeader(r.Header),
		Body:      body,
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.RawQuery,
	}
}

func copyHeader(source map[string][]string) map[string][]string {
	target := make(map[string][]string)
	for key, sourceValues := range source {
		targetValues := make([]string, len(sourceValues))
		copy(targetValues, sourceValues)
		target[key] = targetValues
	}
	return target
}

func bodyToStringAndReader(body io.ReadCloser) (string, io.ReadCloser) {
	b, err := io.ReadAll(body)
	if err != nil {
		return "ERROR WHEN READING BODY", nil
	}
	err = body.Close()
	if err != nil {
		return "ERROR WHEN CLOSING BODY READER", nil
	}
	return string(b), io.NopCloser(bytes.NewReader(b))
}

func (b *Builder) Build() *Server {
	return &Server{server: &http.Server{Addr: b.serverAddr(), Handler: b.Router}}
}

func (b *Builder) serverAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", b.port)
}
