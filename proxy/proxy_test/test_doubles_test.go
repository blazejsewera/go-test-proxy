package proxy_test

import (
	"compress/gzip"
	"github.com/blazejsewera/go-test-proxy/header"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/req"
	"github.com/blazejsewera/go-test-proxy/test/res"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MonitorSpy struct {
	Events []proxy.HTTPEvent
	Errors []error
}

var _ proxy.Monitor = (*MonitorSpy)(nil)

func (m *MonitorSpy) HTTPEvent(event proxy.HTTPEvent) {
	m.Events = append(m.Events, event)
}

func (m *MonitorSpy) Err(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *MonitorSpy) Clear() {
	m.Events = []proxy.HTTPEvent{}
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
		expected := req.ReferenceRequest()
		actual := r
		err := assert.RequestsEqualExcludingPathAndHost(expected, actual)
		if err != nil {
			badRequest(w, err)
			return
		}
		header.CloneToResponseWriter(req.ReferenceResponseHeader(), w)
		must.Succeed(w.Write([]byte(r.URL.Path)))
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

func TestPathEchoServer(t *testing.T) {
	url, closeServer := PathEchoServer()
	defer closeServer()

	requestPath := "/test"
	response := must.Succeed(http.DefaultClient.Do(req.New(url, requestPath)))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body := must.Succeed(io.ReadAll(response.Body))
	assert.Equal(t, requestPath, string(body))
}

// GzipServer constructs a new httptest.Server
// that responds with a gzipped body
// with reference body content.
// See: res.ReferenceBody
func GzipServer() (url string, closeServer func()) {
	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Encoding", "gzip")
		gzipped := gzip.NewWriter(w)
		must.Succeed(gzipped.Write([]byte(res.ReferenceBody())))
		_ = gzipped.Close()
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}
