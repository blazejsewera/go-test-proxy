package proxy

import (
	"bytes"
	"compress/gzip"
	"github.com/blazejsewera/go-test-proxy/header"
	"io"
	"net/http"
)

type ResponseInterceptor struct {
	w          http.ResponseWriter
	bodyBuffer *bytes.Buffer
	statusCode int
	monitor    Monitor
}

var _ http.ResponseWriter = (*ResponseInterceptor)(nil)

func NewResponseWriterInterceptor(w http.ResponseWriter, monitor Monitor) *ResponseInterceptor {
	return &ResponseInterceptor{
		w:          w,
		bodyBuffer: bytes.NewBuffer([]byte{}),
		statusCode: http.StatusOK,
		monitor:    monitor,
	}
}

func (i *ResponseInterceptor) Header() http.Header {
	return i.w.Header()
}

func (i *ResponseInterceptor) Write(body []byte) (int, error) {
	return i.bodyBuffer.Write(body)
}

func (i *ResponseInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.w.WriteHeader(statusCode)
}

func (i *ResponseInterceptor) responseHTTPEvent() HTTPEvent {
	headerCopy := header.Clone(i.w.Header())
	body := i.bodyBufferToString(i.bodyBuffer, headerCopy)
	return HTTPEvent{
		EventType: ResponseEventType,
		Header:    headerCopy,
		Body:      body,
		Status:    i.statusCode,
	}
}

func (i *ResponseInterceptor) bodyBufferToString(bodyBuffer *bytes.Buffer, header map[string][]string) string {
	if gzipped(header) {
		compressed := bytes.NewBuffer(bodyBuffer.Bytes())
		return i.gunzip(compressed)
	} else {
		return bodyBuffer.String()
	}
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
