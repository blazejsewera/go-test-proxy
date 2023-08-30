package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/request"
	"io"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	monitor := new(MonitorSpy)

	t.Run("proxy without any custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		tested := NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitor).
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

		t.Run("monitors forwarded request and response", func(t *testing.T) {
			monitor.Clear()

			requestPath := "/test"
			requestEvent := proxy.HTTPEvent{
				EventType: proxy.RequestEventType,
				Header:    request.ReferenceHeader(),
				Body:      request.ReferenceBody(),
				Method:    request.MethodGet(),
				Path:      requestPath,
			}
			responseEvent := proxy.HTTPEvent{
				EventType: proxy.ResponseEventType,
				Header:    request.ReferenceResponseHeader(),
				Body:      requestPath,
				Status:    http.StatusOK,
			}
			expected := []proxy.HTTPEvent{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(request.New(tested.URL, requestPath)))

			assert.HTTPEventsEqual(t, expected, monitor.Events)
		})
	})

	t.Run("proxy with a custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		customPath := "/customPath"
		customResponseBody := "customResponseBody"
		customHandler := func(w http.ResponseWriter, r *http.Request) {
			must.Succeed(w.Write([]byte(customResponseBody)))
		}

		tested := NewBuilder().
			WithHandlerFunc(customPath, customHandler).
			WithProxyTarget(backendURL).
			WithMonitor(monitor).
			Build()

		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("does not forward a request for a custom path and uses a custom handler instead", func(t *testing.T) {
			response := must.Succeed(client.Do(request.New(tested.URL, customPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, customResponseBody, string(body))
		})

		t.Run("forwards a request with headers to the underlying backend server for a different path", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(request.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := must.Succeed(io.ReadAll(response.Body))
			assert.Equal(t, requestPath, string(body))
		})

		t.Run("monitors request and response handled by custom handler", func(t *testing.T) {
			monitor.Clear()

			requestEvent := proxy.HTTPEvent{
				EventType: proxy.RequestEventType,
				Header:    request.ReferenceHeader(),
				Body:      request.ReferenceBody(),
				Method:    request.MethodGet(),
				Path:      customPath,
			}
			responseEvent := proxy.HTTPEvent{
				EventType: proxy.ResponseEventType,
				Header:    http.Header{},
				Body:      customResponseBody,
				Status:    http.StatusOK,
			}
			expected := []proxy.HTTPEvent{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(request.New(tested.URL, customPath)))

			assert.HTTPEventsEqual(t, expected, monitor.Events)
		})
	})
}
