package interceptor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/ext/brotli"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/header"
	"github.com/blazejsewera/go-test-proxy/test/must"
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
	compressionT := compression(header)
	switch compressionT {
	case gzipCompression:
		compressed := bytes.NewBuffer(i.bodyBuffer.Bytes())
		return i.decompressGzip(compressed)
	case brotliCompression:
		compressed := bytes.NewBuffer(i.bodyBuffer.Bytes())
		return i.decompressBrotli(compressed)
	case noCompression:
		fallthrough
	default:
		return i.bodyBuffer.String()
	}
}

type compressionType int

const (
	noCompression compressionType = iota
	gzipCompression
	brotliCompression
)

func compression(header http.Header) compressionType {
	values, ok := header["Content-Encoding"]
	if !ok {
		return noCompression
	}

	if slices.Contains(values, "gzip") {
		return gzipCompression
	}
	if slices.Contains(values, "br") {
		return brotliCompression
	}

	return noCompression
}

func (i *ResponseInterceptor) decompressGzip(compressed *bytes.Buffer) string {
	decompressed := &bytes.Buffer{}
	gzipReader, err := gzip.NewReader(compressed)
	if err != nil {
		i.monitor.Err(err)
		return ""
	}
	defer must.Close(gzipReader)

	_, err = io.Copy(decompressed, gzipReader)
	if err != nil {
		i.monitor.Err(err)
		return ""
	}

	return decompressed.String()
}

func (i *ResponseInterceptor) decompressBrotli(compressed *bytes.Buffer) string {
	decompressed := &bytes.Buffer{}
	brotliReader := brotli.NewReader(compressed)
	defer must.Close(brotliReader)

	_, err := io.Copy(decompressed, brotliReader)
	if err != nil {
		i.monitor.Err(err)
		return ""
	}

	return decompressed.String()
}
