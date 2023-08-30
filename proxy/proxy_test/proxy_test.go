package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/requests"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	testRequest := func(baseURL string) *http.Request {
		testPath := "/test"
		headers := map[string]string{"X-Test-Header": "Test-Value"}
		return requests.MustMakeNewRequest("GET", baseURL+testPath, "body", headers)
	}

	underlyingBackendServer := func(backendResponse string) (url string, closeBackend func()) {
		badRequest := func(w http.ResponseWriter, err error) {
			w.WriteHeader(http.StatusBadRequest)
			must.Succeed(w.Write([]byte(err.Error())))
		}

		backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expected := testRequest("")
			actual := r
			err := requests.AssertEqualExcludingHost(expected, actual)
			if err != nil {
				badRequest(w, err)
				return
			}
			must.Succeed(w.Write([]byte(backendResponse)))
		})

		backend := httptest.NewServer(backendEndpoint)
		return backend.URL, backend.Close
	}

	t.Run("proxy without any custom handler", func(t *testing.T) {
		backendResponse := "ok"
		backendURL, closeBackend := underlyingBackendServer(backendResponse)
		defer closeBackend()

		tested := NewBuilder().
			WithProxyTarget(backendURL).
			Build()
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards a request with headers to the underlying backend server", func(t *testing.T) {
			request := testRequest(tested.URL)

			response := must.Succeed(client.Do(request))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, backendResponse, string(body))
		})
	})

	t.Run("proxy with a custom handler", func(t *testing.T) {
		backendResponse := "ok"
		backendURL, closeBackend := underlyingBackendServer(backendResponse)
		defer closeBackend()

		customResponse := "customResponse"
		customHandler := func(w http.ResponseWriter, r *http.Request) {
			must.Succeed(w.Write([]byte(customResponse)))
		}

		tested := NewBuilder().
			WithProxyTarget(backendURL).
			WithHandlerFunc("/test", customHandler).
			Build()

		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("does not forward a request for a particular path and uses a custom handler instead", func(t *testing.T) {
			request := testRequest(tested.URL)

			response := must.Succeed(client.Do(request))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, customResponse, string(body))
		})
	})
}
