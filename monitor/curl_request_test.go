package monitor_test

import (
	"bytes"
	"net/http"
	"regexp"
	"testing"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/event"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/test/assert"
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

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cfmt := colorfmt.New(false, stdout, stderr)
		tested := monitor.NewCurlRequestMonitor(target, cfmt)

		tested.HTTPEvent(httpEvent)

		assert.Equal(t, expected, stdout.String())
	})

	t.Run("discards response HTTPEvent", func(t *testing.T) {
		target := ""
		httpEvent := event.HTTP{EventType: event.ResponseEventType}

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cfmt := colorfmt.New(false, stdout, stderr)
		tested := monitor.NewCurlRequestMonitor(target, cfmt)

		tested.HTTPEvent(httpEvent)

		assert.Empty(t, stdout.Bytes())
		assert.Empty(t, stderr.Bytes())
	})
}

func normalizedLine(s string) string {
	spaces := regexp.MustCompile(`\s+`)
	return spaces.ReplaceAllString(s, " ") + "\n"
}
