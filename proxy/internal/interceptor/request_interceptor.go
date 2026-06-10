package interceptor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/header"
)

type Request struct {
	request *http.Request
	monitor monitor.Monitor
}

func NewRequestMonitor(r *http.Request, monitor monitor.Monitor) *Request {
	return &Request{r, monitor}
}

func (i *Request) MonitorRequest() {
	i.monitor.HTTPEvent(i.requestHTTPEvent())
}

func (i *Request) requestHTTPEvent() event.HTTP {
	h := http.Header{}
	header.Copy(h, i.request.Header)

	body, bodyReader := i.bodyToStringAndReader(i.request.Body)
	i.request.Body = bodyReader
	return event.HTTP{
		EventType: event.RequestEventType,
		Header:    h,
		Body:      body,
		Method:    i.request.Method,
		Path:      i.request.URL.Path,
		Query:     i.request.URL.RawQuery,
	}
}

func (i *Request) bodyToStringAndReader(body io.ReadCloser) (string, io.ReadCloser) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		i.monitor.Err(fmt.Errorf("read request body: %s", err))
		return "", nil
	}
	err = body.Close()
	if err != nil {
		i.monitor.Err(fmt.Errorf("close request body: %s", err))
		return "", nil
	}
	return string(bodyBytes), io.NopCloser(bytes.NewReader(bodyBytes))
}
