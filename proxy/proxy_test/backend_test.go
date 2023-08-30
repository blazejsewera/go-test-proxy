package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/header"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/request"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// PathEchoServer constructs a new httptest.Server
// that responds with the Path of a Request it received
func PathEchoServer() (url string, closeServer func()) {
	badRequest := func(w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusBadRequest)
		must.Succeed(w.Write([]byte(err.Error())))
	}

	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := request.ReferenceRequest()
		actual := r
		err := assert.RequestsEqualExcludingPathAndHost(expected, actual)
		if err != nil {
			badRequest(w, err)
			return
		}
		header.CloneToResponseWriter(request.ReferenceResponseHeader(), w)
		must.Succeed(w.Write([]byte(r.URL.Path)))
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

func TestPathEchoServer(t *testing.T) {
	url, closeServer := PathEchoServer()
	defer closeServer()

	requestPath := "/test"
	response := must.Succeed(http.DefaultClient.Do(request.New(url, requestPath)))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body := must.Succeed(io.ReadAll(response.Body))
	assert.Equal(t, requestPath, string(body))
}
