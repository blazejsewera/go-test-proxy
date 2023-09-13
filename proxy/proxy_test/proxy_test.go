package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/header"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/req"
	"github.com/blazejsewera/go-test-proxy/test/res"
	"io"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	monitor := new(MonitorSpy)

	t.Run("proxy without any custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		tested := BuildTestServer(proxy.NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitor))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards a request with headers to the underlying backend server", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(req.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, requestPath, string(body))
			assert.HeaderContainsExpected(t, req.ReferenceResponseHeader(), response.Header)
		})

		t.Run("monitors forwarded request and response", func(t *testing.T) {
			monitor.Clear()

			requestPath := "/test"
			requestEvent := proxy.HTTPEvent{
				EventType: proxy.RequestEventType,
				Header:    req.ReferenceHeader(),
				Body:      req.ReferenceBody(),
				Method:    req.MethodGet(),
				Path:      requestPath,
			}
			responseEvent := proxy.HTTPEvent{
				EventType: proxy.ResponseEventType,
				Header:    req.ReferenceResponseHeader(),
				Body:      requestPath,
				Status:    http.StatusOK,
			}
			expected := []proxy.HTTPEvent{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(req.New(tested.URL, requestPath)))

			assert.HTTPEventsEqual(t, expected, monitor.Events)
			assert.Empty(t, monitor.Errors)
		})
	})

	t.Run("proxy with a custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		customPath := "/customPath"
		customResponseBody := "customResponseBody"
		customResponseHeader := http.Header{"X-Custom-Header": []string{"Custom-Value"}}
		customHandler := func(w http.ResponseWriter, r *http.Request) {
			header.CloneToResponseWriter(customResponseHeader, w)
			must.Succeed(w.Write([]byte(customResponseBody)))
		}

		tested := BuildTestServer(proxy.NewBuilder().
			WithHandlerFunc(customPath, customHandler).
			WithProxyTarget(backendURL).
			WithMonitor(monitor))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("does not forward a request for a custom path and uses a custom handler instead", func(t *testing.T) {
			response := must.Succeed(client.Do(req.New(tested.URL, customPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, customResponseBody, string(body))
			assert.HeaderContainsExpected(t, customResponseHeader, response.Header)
		})

		t.Run("forwards a request with headers to the underlying backend server for a different path", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(req.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, requestPath, string(body))
			assert.HeaderContainsExpected(t, req.ReferenceResponseHeader(), response.Header)
		})

		t.Run("monitors request and response handled by custom handler", func(t *testing.T) {
			monitor.Clear()

			requestEvent := proxy.HTTPEvent{
				EventType: proxy.RequestEventType,
				Header:    req.ReferenceHeader(),
				Body:      req.ReferenceBody(),
				Method:    req.MethodGet(),
				Path:      customPath,
			}
			responseEvent := proxy.HTTPEvent{
				EventType: proxy.ResponseEventType,
				Header:    http.Header{},
				Body:      customResponseBody,
				Status:    http.StatusOK,
			}
			expected := []proxy.HTTPEvent{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(req.New(tested.URL, customPath)))

			assert.HTTPEventsEqual(t, expected, monitor.Events)
			assert.Empty(t, monitor.Errors)
		})
	})

	t.Run("proxy forwarding gzip-compressed payload", func(t *testing.T) {
		monitor.Clear()

		backendURL, closeBackend := GzipServer()
		defer closeBackend()

		tested := BuildTestServer(proxy.NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitor))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards it unchanged but monitors it in plain text", func(t *testing.T) {
			response := must.Succeed(client.Do(req.New(tested.URL, "/")))

			// body gets automatically unzipped by http.Client
			// based on 'Content-Encoding: gzip' header value
			actual := readBody(response.Body)
			assert.Equal(t, res.ReferenceBody(), actual)

			responseEventBody := monitor.Events[1].Body
			assert.Equal(t, res.ReferenceBody(), responseEventBody)
			assert.Empty(t, monitor.Errors)
		})
	})
}

func readBody(r io.Reader) string {
	b := must.Succeed(io.ReadAll(r))
	return string(b)
}
