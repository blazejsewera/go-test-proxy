package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy/proxytest"
	"github.com/blazejsewera/go-test-proxy/requests"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const expectedResponse = "ok"

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler functions forwards a request with headers to the underlying endpoint",
		func(t *testing.T) {
			backendURL, closeBackend := backendServer()
			defer closeBackend()

			var tested = proxytest.Builder().
				WithTarget(backendURL).
				Build()
			tested.Start()
			defer tested.Close()

			var client = tested.Client()
			request := requestStub(tested.URL)
			response := must.Succeed(client.Do(request))

			body := must.Succeed(io.ReadAll(response.Body))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, expectedResponse, string(body))
		})
}

func backendServer() (url string, closeBackend func()) {
	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := requestStub("")
		actual := r
		err := requests.AssertEqualExcludingHost(expected, actual)
		if err != nil {
			badRequest(w, err)
			return
		}
		_, err = w.Write([]byte(expectedResponse))
		if err != nil {
			panic("write response in backend endpoint: " + err.Error())
		}
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		panic("write response in backend: " + err.Error())
	}
	return
}

func requestStub(baseURL string) *http.Request {
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return requests.MustMakeNewRequest("GET", baseURL+"/test", "body", headers)
}
