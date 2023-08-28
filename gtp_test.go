package gtp_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy/proxytest"
	"github.com/blazejsewera/go-test-proxy/requests"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const expectedResponse = "ok"

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler functions forwards a request with headers to the underlying endpoint",
		func(t *testing.T) {
			backendURL, closeBackend := backendServer(t)
			defer closeBackend()

			var tested = proxytest.Builder().
				WithTarget(backendURL).
				Build()

			tested.Start()
			defer tested.Close()

			var client = tested.Client()
			request := requestStub(t, tested.URL)
			responseStruct, err := client.Do(request)
			if err != nil {
				t.Fatal("send request: ", err)
			}

			response, err := io.ReadAll(responseStruct.Body)
			if err != nil {
				t.Fatal("read request: ", err)
			}

			if responseStruct.StatusCode != http.StatusOK {
				t.Fatal("response status code not OK: ", responseStruct.StatusCode)
			}
			if string(response) != expectedResponse {
				t.Fatalf("response not equal: %s, %s", expectedResponse, string(response))
			}
		})
}

func backendServer(t testing.TB) (url string, closeServer func()) {
	backendEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := requestStub(t, "http://excluded-from-test")
		actual := r
		err := requests.AssertEqualExcludingHost(expected, actual)
		if err != nil {
			badRequest(t, w, err)
			return
		}
		_, err = w.Write([]byte(expectedResponse))
		if err != nil {
			t.Fatal("write response in backend endpoint: ", err)
		}
	})

	backend := httptest.NewServer(backendEndpoint)
	return backend.URL, backend.Close
}

func badRequest(t testing.TB, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		t.Fatal("write response in backend: ", err)
	}
	return
}

func requestStub(t testing.TB, baseURL string) *http.Request {
	t.Helper()
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return requests.MustMakeNewRequest(t, "GET", baseURL+"/test", "body", headers)
}
