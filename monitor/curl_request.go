package monitor

import (
	"fmt"
	event2 "github.com/blazejsewera/go-test-proxy/event"
	"io"
	"os"
	"strconv"
)

type curlRequest struct {
	target string
	output io.Writer
}

func NewCurlRequestMonitor(target string) Monitor {
	return NewCurlRequestMonitorW(target, os.Stdout)
}

func NewCurlRequestMonitorW(target string, output io.Writer) Monitor {
	return &curlRequest{target: target, output: output}
}

func (c *curlRequest) HTTPEvent(event event2.HTTP) {
	switch event.EventType {
	case event2.RequestEventType:
		c.writeCurlRequest(event)
	case event2.ResponseEventType:
		return
	default:
		return
	}
}

func (c *curlRequest) writeCurlRequest(event event2.HTTP) {
	result := fmt.Sprintf("curl -X %s%s%s %s\n",
		event.Method,
		headerToCurl(event.Header),
		bodyToCurl(event.Body),
		urlPartsToCurl(c.target, event.Path, event.Query))
	_, err := c.output.Write([]byte(result))
	if err != nil {
		c.Err(err)
		return
	}
}

func headerToCurl(header map[string][]string) string {
	result := ""
	for key, values := range header {
		for _, value := range values {
			result += fmt.Sprintf(" -H \"%s: %s\"", key, value)
		}
	}
	return result
}

func bodyToCurl(body string) string {
	if body == "" {
		return ""
	}
	return fmt.Sprintf(" -d %s", strconv.Quote(body))
}

func urlPartsToCurl(host string, path string, query string) string {
	if query == "" {
		return host + path
	}
	return fmt.Sprintf("%s%s?%s", host, path, query)
}

func (c *curlRequest) Err(error) {}
