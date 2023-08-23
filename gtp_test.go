package gtp_test

import (
	"github.com/blazejsewera/go-test-proxy/requests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler functions forwards a request with headers to the underlying endpoint",
		func(t *testing.T) {
			// given
			backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expected := requestStub(t, "http://excluded-from-test")
				actual := r
				requests.AssertEqualExcludingHost(t, expected, actual)
			})
			backend := httptest.NewServer(backendEndpoint)
			backend.Start()

			var tested *http.Server = gtptest.Builder().
				WithTarget(backend.URL).
				Build()

			backend.Client()
			var client *http.Client = tested.Client()

			tested.Start()

			// when
			var baseURL string = tested.URL
			headers := map[string]string{"X-Test-Header": "Test-Value"}
			request := requestStub(t, baseURL)
			client.Do(request)
		})
}

func requestStub(t testing.TB, baseURL string) *http.Request {
	t.Helper()
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return requests.MustMakeNewRequest(t, "GET", baseURL+"/test", "body", headers)
}
