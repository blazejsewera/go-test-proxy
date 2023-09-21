package monitor_test

import (
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestCurlRequestMonitor(t *testing.T) {
	t.Run("generates a valid curl cmd given request HTTPEvent and target URL", func(t *testing.T) {
		target := "http://example.com"
		httpEvent := event.HTTP{
			EventType: event.RequestEventType,
			Header: map[string][]string{
				"Accept":       {"application/json"},
				"Content-Type": {"application/json"},
			},
			Body:   `{"bodyKey":"bodyValue"}`,
			Method: http.MethodPost,
			Path:   "/path",
			Query:  "queryKey=queryValue",
		}

		expected := normalizedLine(
			`curl -X POST
			      -H "Accept: application/json"
			      -H "Content-Type: application/json"
			      -d "{\"bodyKey\":\"bodyValue\"}"
    		      http://example.com/path?queryKey=queryValue`)

		buffer := &strings.Builder{}
		tested := monitor.NewCurlRequestMonitorW(target, buffer)

		tested.HTTPEvent(httpEvent)

		assert.Equal(t, expected, buffer.String())
	})

	t.Run("discards response HTTPEvent", func(t *testing.T) {
		target := ""
		httpEvent := event.HTTP{EventType: event.ResponseEventType}

		buffer := &strings.Builder{}
		tested := monitor.NewCurlRequestMonitorW(target, buffer)

		tested.HTTPEvent(httpEvent)

		assert.Equal(t, "", buffer.String())
	})
}

func normalizedLine(s string) string {
	spaces := regexp.MustCompile(`\s+`)
	return spaces.ReplaceAllString(s, " ") + "\n"
}
