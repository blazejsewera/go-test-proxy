package proxy

import (
	"bytes"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/header"
	"io"
	"net/http"
)

type requestInterceptor struct {
	request *http.Request
	monitor Monitor
}

func newRequestInterceptor(r *http.Request, monitor Monitor) *requestInterceptor {
	return &requestInterceptor{r, monitor}
}

func (i *requestInterceptor) monitorRequest() {
	i.monitor.HTTPEvent(i.requestHTTPEvent())
}

func (i *requestInterceptor) requestHTTPEvent() HTTPEvent {
	h := http.Header{}
	header.Copy(h, i.request.Header)

	body, bodyReader := i.bodyToStringAndReader(i.request.Body)
	i.request.Body = bodyReader
	return HTTPEvent{
		EventType: RequestEventType,
		Header:    h,
		Body:      body,
		Method:    i.request.Method,
		Path:      i.request.URL.Path,
		Query:     i.request.URL.RawQuery,
	}
}

func (i *requestInterceptor) bodyToStringAndReader(body io.ReadCloser) (string, io.ReadCloser) {
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
