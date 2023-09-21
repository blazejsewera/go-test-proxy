package interceptor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/header"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"io"
	"net/http"
	"slices"
)

type ResponseInterceptor struct {
	responseWriter http.ResponseWriter
	statusCode     int
	bodyBuffer     bytes.Buffer
	monitor        monitor.Monitor
}

var _ http.ResponseWriter = (*ResponseInterceptor)(nil)

func NewResponseInterceptor(w http.ResponseWriter, monitor monitor.Monitor) *ResponseInterceptor {
	return &ResponseInterceptor{
		responseWriter: w,
		statusCode:     http.StatusOK,
		bodyBuffer:     bytes.Buffer{},
		monitor:        monitor,
	}
}

func (i *ResponseInterceptor) Header() http.Header {
	return i.responseWriter.Header()
}

func (i *ResponseInterceptor) Write(body []byte) (int, error) {
	return i.bodyBuffer.Write(body)
}

func (i *ResponseInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
}

func (i *ResponseInterceptor) MonitorAndForwardResponse() {
	i.monitor.HTTPEvent(i.responseHTTPEvent())

	i.responseWriter.WriteHeader(i.statusCode)
	_, err := io.Copy(i.responseWriter, &i.bodyBuffer)
	if err != nil {
		i.monitor.Err(fmt.Errorf("copy interceptor buffer to response writer: %s", err))
	}
}

func (i *ResponseInterceptor) responseHTTPEvent() event.HTTP {
	h := http.Header{}
	header.Copy(h, i.responseWriter.Header())
	body := i.bodyBufferToString(h)
	return event.HTTP{
		EventType: event.ResponseEventType,
		Header:    h,
		Body:      body,
		Status:    i.statusCode,
	}
}

func (i *ResponseInterceptor) bodyBufferToString(header map[string][]string) string {
	if gzipped(header) {
		compressed := bytes.NewBuffer(i.bodyBuffer.Bytes())
		return i.gunzip(compressed)
	} else {
		return i.bodyBuffer.String()
	}
}

func gzipped(header map[string][]string) bool {
	result := false
	values, ok := header["Content-Encoding"]
	if ok {
		result = slices.Contains(values, "gzip")
	}
	return result
}

func (i *ResponseInterceptor) gunzip(compressed *bytes.Buffer) string {
	decompressed := &bytes.Buffer{}
	gzipReader, err := gzip.NewReader(compressed)
	if err != nil {
		i.monitor.Err(err)
		return ""
	}

	_, err = io.Copy(decompressed, gzipReader)
	if err != nil {
		i.monitor.Err(err)
		return ""
	}

	err = gzipReader.Close()
	if err != nil {
		i.monitor.Err(err)
		return ""
	}

	return decompressed.String()
}
