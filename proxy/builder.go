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
		port:    8080,
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

		r.Host = targetURL.Host
		r.RequestURI = ""
		r.URL = targetURL

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			return
		}
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return
		}
		_, err = w.Write(bodyBytes)
		if err != nil {
			return
		}
		header.CloneToResponseWriter(response.Header, w)
	}

	return b.WithHandlerFunc("/", proxyHandler)
}

type ResponseWriterInterceptor struct {
	w          http.ResponseWriter
	bodyBuffer *bytes.Buffer
	statusCode int
}

var _ http.ResponseWriter = (*ResponseWriterInterceptor)(nil)

func NewResponseWriterInterceptor(w http.ResponseWriter) *ResponseWriterInterceptor {
	return &ResponseWriterInterceptor{
		w:          w,
		bodyBuffer: bytes.NewBuffer([]byte{}),
		statusCode: http.StatusOK,
	}
}

func (i *ResponseWriterInterceptor) Header() http.Header {
	return i.w.Header()
}

func (i *ResponseWriterInterceptor) Write(body []byte) (int, error) {
	return i.bodyBuffer.Write(body)
}

func (i *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.w.WriteHeader(statusCode)
}

func (b *Builder) WithHandlerFunc(pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) *Builder {
	wrapperFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.Monitor.HTTPEvent(requestHTTPEvent(r))
		interceptor := NewResponseWriterInterceptor(w)
		handlerFunc(interceptor, r)
		b.Monitor.HTTPEvent(responseHTTPEvent(interceptor))

		_, err := io.Copy(w, interceptor.bodyBuffer)
		if err != nil {
			return
		}
	})
	b.Router.Handle(pattern, wrapperFunc)
	return b
}

func responseHTTPEvent(interceptor *ResponseWriterInterceptor) HTTPEvent {
	body := interceptor.bodyBuffer.String()
	return HTTPEvent{
		EventType: ResponseEventType,
		Header:    header.Clone(interceptor.w.Header()),
		Body:      body,
		Status:    interceptor.statusCode,
	}
}

func requestHTTPEvent(r *http.Request) HTTPEvent {
	body, bodyReader := bodyToStringAndReader(r.Body)
	r.Body = bodyReader
	return HTTPEvent{
		EventType: RequestEventType,
		Header:    header.Clone(r.Header),
		Body:      body,
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.RawQuery,
	}
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
