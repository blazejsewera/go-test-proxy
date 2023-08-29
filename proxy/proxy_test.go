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

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler"+
		" forwards a request with headers to the underlying backend server",
		func(t *testing.T) {
			backendResponse := "ok"
			backendURL, closeBackend := backendServer(backendResponse)
			defer closeBackend()

			tested := proxytest.Builder().
				WithTarget(backendURL).
				Build()
			tested.Start()
			defer tested.Close()

			client := tested.Client()
			request := testRequest(tested.URL)
			response := must.Succeed(client.Do(request))

			body := must.Succeed(io.ReadAll(response.Body))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, backendResponse, string(body))
		})

	t.Run("proxy with a custom handler"+
		" does not forward a request for a particular path"+
		" and uses a custom handler instead",
		func(t *testing.T) {
			backendResponse := "ok"
			backendURL, closeBackend := backendServer(backendResponse)
			defer closeBackend()

			customResponse := "customResponse"
			customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				must.Succeed(w.Write([]byte(customResponse)))
			})

			tested := proxytest.Builder().
				WithTarget(backendURL).
				WithHandler("/test", customHandler).
				Build()

			tested.Start()
			defer tested.Close()

			client := tested.Client()
			request := testRequest(tested.URL)
			response := must.Succeed(client.Do(request))

			body := must.Succeed(io.ReadAll(response.Body))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, customResponse, string(body))
		})
}

func backendServer(backendResponse string) (url string, closeBackend func()) {
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

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	must.Succeed(w.Write([]byte(err.Error())))
}

func testRequest(baseURL string) *http.Request {
	testPath := "/test"
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return requests.MustMakeNewRequest("GET", baseURL+testPath, "body", headers)
}
