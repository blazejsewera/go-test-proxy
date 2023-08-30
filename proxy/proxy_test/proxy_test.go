package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/request"
	"io"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	t.Run("proxy without any custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		tested := NewBuilder().
			WithProxyTarget(backendURL).
			Build()
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards a request with headers to the underlying backend server", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(request.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, requestPath, string(body))
		})
	})

	t.Run("proxy with a custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		customPath := "/customPath"
		customResponse := "customResponse"
		customHandler := func(w http.ResponseWriter, r *http.Request) {
			must.Succeed(w.Write([]byte(customResponse)))
		}

		tested := NewBuilder().
			WithHandlerFunc(customPath, customHandler).
			WithProxyTarget(backendURL).
			Build()

		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("does not forward a request for a custom path and uses a custom handler instead", func(t *testing.T) {
			response := must.Succeed(client.Do(request.New(tested.URL, customPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, customResponse, string(body))
		})

		t.Run("forwards a request with headers to the underlying backend server for a different path", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(request.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, requestPath, string(body))
		})
	})
}
