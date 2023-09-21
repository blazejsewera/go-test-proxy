package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/header"
	"io"
	"net/http"
	"slices"
)

type responseInterceptor struct {
	responseWriter http.ResponseWriter
	statusCode     int
	bodyBuffer     bytes.Buffer
	monitor        Monitor
}

var _ http.ResponseWriter = (*responseInterceptor)(nil)

func newResponseInterceptor(w http.ResponseWriter, monitor Monitor) *responseInterceptor {
	return &responseInterceptor{
		responseWriter: w,
		statusCode:     http.StatusOK,
		bodyBuffer:     bytes.Buffer{},
		monitor:        monitor,
	}
}

func (i *responseInterceptor) Header() http.Header {
	return i.responseWriter.Header()
}

func (i *responseInterceptor) Write(body []byte) (int, error) {
	return i.bodyBuffer.Write(body)
}

func (i *responseInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
}

func (i *responseInterceptor) monitorAndForwardResponse() {
	i.monitor.HTTPEvent(i.responseHTTPEvent())

	i.responseWriter.WriteHeader(i.statusCode)
	_, err := io.Copy(i.responseWriter, &i.bodyBuffer)
	if err != nil {
		i.monitor.Err(fmt.Errorf("copy interceptor buffer to response writer: %s", err))
	}
}

func (i *responseInterceptor) responseHTTPEvent() HTTPEvent {
	h := http.Header{}
	header.Copy(h, i.responseWriter.Header())
	body := i.bodyBufferToString(h)
	return HTTPEvent{
		EventType: ResponseEventType,
		Header:    h,
		Body:      body,
		Status:    i.statusCode,
	}
}

func (i *responseInterceptor) bodyBufferToString(header map[string][]string) string {
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

func (i *responseInterceptor) gunzip(compressed *bytes.Buffer) string {
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
