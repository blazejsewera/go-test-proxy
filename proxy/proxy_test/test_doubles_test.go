package proxy_test

import (
	"compress/gzip"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/monitor/event"
	"github.com/blazejsewera/go-test-proxy/proxy/header"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/testrequest"
	"github.com/blazejsewera/go-test-proxy/test/testresponse"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MonitorSpy struct {
	Events []event.HTTP
	Errors []error
}

var _ monitor.Monitor = (*MonitorSpy)(nil)

func (m *MonitorSpy) HTTPEvent(event event.HTTP) {
	m.Events = append(m.Events, event)
}

func (m *MonitorSpy) Err(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *MonitorSpy) Clear() {
	m.Events = []event.HTTP{}
	m.Errors = []error{}
}

// PathEchoServer constructs a new httptest.Server
// that responds with the Path of a Request it received
func PathEchoServer() (url string, closeServer func()) {
	badRequest := func(w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusBadRequest)
		must.Succeed(w.Write([]byte(err.Error())))
	}

	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := testrequest.ReferenceRequest()
		actual := r
		err := assert.RequestsEqualExcludingPathAndHost(expected, actual)
		if err != nil {
			badRequest(w, err)
			return
		}
		header.Copy(w.Header(), testrequest.ReferenceResponseHeader())
		must.Succeed(w.Write([]byte(r.URL.Path)))
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

func TestPathEchoServer(t *testing.T) {
	url, closeServer := PathEchoServer()
	defer closeServer()

	requestPath := "/test"
	response := must.Succeed(http.DefaultClient.Do(testrequest.New(url, requestPath)))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body := must.Succeed(io.ReadAll(response.Body))
	assert.Equal(t, requestPath, string(body))
}

// NotFoundServer constructs a new httptest.Server
// that always responds with 404 Not Found.
func NotFoundServer() (url string, closeServer func()) {
	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

// GzipServer constructs a new httptest.Server
// that responds with a gzipped body
// with reference body content.
// See: testresponse.ReferenceBody
func GzipServer() (url string, closeServer func()) {
	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Encoding", "gzip")
		gzipped := gzip.NewWriter(w)
		must.Succeed(gzipped.Write([]byte(testresponse.ReferenceBody())))
		_ = gzipped.Close()
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}
