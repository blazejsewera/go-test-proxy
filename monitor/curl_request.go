package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/monitor/event"
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

func (c *curlRequest) HTTPEvent(e event.HTTP) {
	switch e.EventType {
	case event.RequestEventType:
		c.writeCurlRequest(e)
	case event.ResponseEventType:
		return
	default:
		return
	}
}

func (c *curlRequest) writeCurlRequest(e event.HTTP) {
	result := fmt.Sprintf("curl -X %s%s%s %s\n",
		e.Method,
		headerToCurl(e.Header),
		bodyToCurl(e.Body),
		urlPartsToCurl(c.target, e.Path, e.Query))
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
