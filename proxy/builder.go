package proxy

import (
	"bytes"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/header"
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
		port:    8000,
		Monitor: DefaultMonitor{},
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

		r.RequestURI = ""
		r.Host = targetURL.Host
		r.URL = targetURL

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			b.Monitor.Err(fmt.Errorf("client request to target: %s", err))
			return
		}
		w.WriteHeader(response.StatusCode)
		_, err = io.Copy(w, response.Body)
		if err != nil {
			b.Monitor.Err(fmt.Errorf("write response: %s", err))
			return
		}
		header.Clone(w.Header(), response.Header)
	}

	return b.WithHandlerFunc("/", proxyHandler)
}

func (b *Builder) WithHandlerFunc(pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.Monitor.HTTPEvent(b.requestHTTPEvent(r))
		interceptor := newResponseInterceptor(w, b.Monitor)
		handlerFunc(interceptor, r)
		b.Monitor.HTTPEvent(interceptor.responseHTTPEvent())

		w.WriteHeader(interceptor.statusCode)
		_, err := io.Copy(w, &interceptor.bodyBuffer)
		if err != nil {
			b.Monitor.Err(fmt.Errorf("copy interceptor buffer to response writer: %s", err))
			return
		}
	})
	b.Router.Handle(pattern, wrapperFunc)
	return b
}

func (b *Builder) requestHTTPEvent(r *http.Request) HTTPEvent {
	h := http.Header{}
	header.Clone(h, r.Header)

	body, bodyReader := b.bodyToStringAndReader(r.Body)
	r.Body = bodyReader
	return HTTPEvent{
		EventType: RequestEventType,
		Header:    h,
		Body:      body,
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.RawQuery,
	}
}

func (b *Builder) bodyToStringAndReader(body io.ReadCloser) (string, io.ReadCloser) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		b.Monitor.Err(fmt.Errorf("read request body: %s", err))
		return "", nil
	}
	err = body.Close()
	if err != nil {
		b.Monitor.Err(fmt.Errorf("close request body: %s", err))
		return "", nil
	}
	return string(bodyBytes), io.NopCloser(bytes.NewReader(bodyBytes))
}

func (b *Builder) Build() *http.Server {
	return &http.Server{Addr: b.serverAddr(), Handler: b.Router}
}

func (b *Builder) serverAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", b.port)
}
