package proxy

import (
	"bytes"
	"compress/gzip"
	"github.com/blazejsewera/go-test-proxy/header"
	"io"
	"net/http"
	"slices"
)

type responseInterceptor struct {
	w          http.ResponseWriter
	bodyBuffer bytes.Buffer
	statusCode int
	monitor    Monitor
}

var _ http.ResponseWriter = (*responseInterceptor)(nil)

func newResponseInterceptor(w http.ResponseWriter, monitor Monitor) *responseInterceptor {
	return &responseInterceptor{
		w:          w,
		bodyBuffer: bytes.Buffer{},
		statusCode: http.StatusOK,
		monitor:    monitor,
	}
}

func (i *responseInterceptor) Header() http.Header {
	return i.w.Header()
}

func (i *responseInterceptor) Write(body []byte) (int, error) {
	return i.bodyBuffer.Write(body)
}

func (i *responseInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.w.WriteHeader(statusCode)
}

func (i *responseInterceptor) responseHTTPEvent() HTTPEvent {
	headerCopy := header.Clone(i.w.Header())
	body := i.bodyBufferToString(headerCopy)
	return HTTPEvent{
		EventType: ResponseEventType,
		Header:    headerCopy,
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
