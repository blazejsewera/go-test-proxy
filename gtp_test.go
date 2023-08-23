package gtp_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy/proxytest"
	"github.com/blazejsewera/go-test-proxy/requests"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler functions forwards a request with headers to the underlying endpoint",
		func(t *testing.T) {
			const expectedResponse = "ok"
			done := make(chan struct{})
			backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expected := requestStub(t, "http://excluded-from-test")
				actual := r
				requests.AssertEqualExcludingHost(t, expected, actual)
				_, err := w.Write([]byte(expectedResponse))
				if err != nil {
					t.Fatalf("write response in backend endpoint: %s", err)
				}
				done <- struct{}{}
			})
			backend := httptest.NewServer(backendEndpoint)
			defer backend.Close()

			var tested = proxytest.Builder().
				WithTarget(backend.URL).
				Build()

			backend.Client()

			tested.Start()
			defer tested.Close()

			var client = tested.Client()
			request := requestStub(t, tested.URL)

			responseStruct, err1 := client.Do(request)
			response, err2 := io.ReadAll(responseStruct.Body)

			if err1 != nil {
				t.Fatalf("request to tested: %s", err1)
			}
			if err2 != nil {
				t.Fatalf("reading response: %s", err2)
			}

			if string(response) != expectedResponse {
				t.Fatalf("response not equal: %s, %s", expectedResponse, string(response))
			}

			select {
			case <-time.After(300 * time.Millisecond):
				t.Fatalf("timeout reached")
			case <-done:
			}
		})
}

func requestStub(t testing.TB, baseURL string) *http.Request {
	t.Helper()
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return requests.MustMakeNewRequest(t, "GET", baseURL+"/test", "body", headers)
}
