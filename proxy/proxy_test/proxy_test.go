package proxy_test

import (
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/proxy/header"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"github.com/blazejsewera/go-test-proxy/test/testrequest"
	"github.com/blazejsewera/go-test-proxy/test/testresponse"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	monitorSpy := new(MonitorSpy)

	t.Run("proxy without any custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		tested := buildTestServer(proxy.NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitorSpy))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards a request with headers to the underlying backend server", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(testrequest.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, requestPath, body)
			assert.HeaderContainsExpected(t, testrequest.ReferenceResponseHeader(), response.Header)
		})

		t.Run("monitors forwarded request and response", func(t *testing.T) {
			monitorSpy.Clear()

			requestPath := "/test"
			requestEvent := event.HTTP{
				EventType: event.RequestEventType,
				Header:    testrequest.ReferenceHeader(),
				Body:      testrequest.ReferenceBody(),
				Method:    testrequest.MethodGet(),
				Path:      requestPath,
			}
			responseEvent := event.HTTP{
				EventType: event.ResponseEventType,
				Header:    testrequest.ReferenceResponseHeader(),
				Body:      requestPath,
				Status:    http.StatusOK,
			}
			expected := []event.HTTP{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(testrequest.New(tested.URL, requestPath)))

			assert.HTTPEventListEqual(t, expected, monitorSpy.Events)
			assert.Empty(t, monitorSpy.Errors)
		})
	})

	t.Run("proxy with a target that sends 404 Not Found status", func(t *testing.T) {
		backendURL, closeBackend := NotFoundServer()
		defer closeBackend()

		tested := buildTestServer(proxy.NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitorSpy))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards and monitors 404 status code", func(t *testing.T) {
			monitorSpy.Clear()

			requestPath := "/test"
			expectedResponseEvent := event.HTTP{
				EventType: event.ResponseEventType,
				Status:    http.StatusNotFound,
			}

			response := must.Succeed(client.Do(testrequest.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			actualResponseEvent := monitorSpy.Events[1]
			assert.HTTPEventsEqual(t, expectedResponseEvent, actualResponseEvent)
			assert.Empty(t, monitorSpy.Errors)
		})
	})

	t.Run("proxy with a custom handler", func(t *testing.T) {
		backendURL, closeBackend := PathEchoServer()
		defer closeBackend()

		customPath := "/customPath"
		customResponseBody := "customResponseBody"
		customResponseHeader := http.Header{"X-Custom-Header": []string{"Custom-Value"}}
		customHandler := func(w http.ResponseWriter, r *http.Request) {
			header.Copy(w.Header(), customResponseHeader)
			must.Succeed(w.Write([]byte(customResponseBody)))
		}

		tested := buildTestServer(proxy.NewBuilder().
			WithHandlerFunc(customPath, customHandler).
			WithProxyTarget(backendURL).
			WithMonitor(monitorSpy))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("does not forward a request for a custom path and uses a custom handler instead", func(t *testing.T) {
			response := must.Succeed(client.Do(testrequest.New(tested.URL, customPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, customResponseBody, body)
			assert.HeaderContainsExpected(t, customResponseHeader, response.Header)
		})

		t.Run("forwards a request with headers to the underlying backend server for a different path", func(t *testing.T) {
			requestPath := "/test"
			response := must.Succeed(client.Do(testrequest.New(tested.URL, requestPath)))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			body := readBody(response.Body)
			assert.Equal(t, requestPath, body)
			assert.HeaderContainsExpected(t, testrequest.ReferenceResponseHeader(), response.Header)
		})

		t.Run("monitors request and response handled by custom handler", func(t *testing.T) {
			monitorSpy.Clear()

			requestEvent := event.HTTP{
				EventType: event.RequestEventType,
				Header:    testrequest.ReferenceHeader(),
				Body:      testrequest.ReferenceBody(),
				Method:    testrequest.MethodGet(),
				Path:      customPath,
			}
			responseEvent := event.HTTP{
				EventType: event.ResponseEventType,
				Header:    http.Header{},
				Body:      customResponseBody,
				Status:    http.StatusOK,
			}
			expected := []event.HTTP{requestEvent, responseEvent}

			_ = must.Succeed(client.Do(testrequest.New(tested.URL, customPath)))

			assert.HTTPEventListEqual(t, expected, monitorSpy.Events)
			assert.Empty(t, monitorSpy.Errors)
		})
	})

	t.Run("proxy forwarding gzip-compressed payload", func(t *testing.T) {
		monitorSpy.Clear()

		backendURL, closeBackend := GzipServer()
		defer closeBackend()

		tested := buildTestServer(proxy.NewBuilder().
			WithProxyTarget(backendURL).
			WithMonitor(monitorSpy))
		tested.Start()
		defer tested.Close()

		client := tested.Client()

		t.Run("forwards it unchanged but monitors it in plain text", func(t *testing.T) {
			response := must.Succeed(client.Do(testrequest.New(tested.URL, "/")))

			// body gets automatically unzipped by http.Client
			// based on 'Content-Encoding: gzip' header value
			actual := readBody(response.Body)
			assert.Equal(t, testresponse.ReferenceBody(), actual)

			responseEventBody := monitorSpy.Events[1].Body
			assert.Equal(t, testresponse.ReferenceBody(), responseEventBody)
			assert.Empty(t, monitorSpy.Errors)
		})
	})
}

func buildTestServer(builder *proxy.Builder) *httptest.Server {
	return httptest.NewUnstartedServer(builder.Router)
}

func readBody(r io.Reader) string {
	b := must.Succeed(io.ReadAll(r))
	return string(b)
}
