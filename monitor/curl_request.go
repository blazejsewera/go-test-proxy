package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"io"
	"os"
	"strconv"
)

type curlRequest struct {
	target string
	output io.Writer
}

func NewCurlRequestMonitor(target string) proxy.Monitor {
	return NewCurlRequestMonitorW(target, os.Stdout)
}

func NewCurlRequestMonitorW(target string, output io.Writer) proxy.Monitor {
	return &curlRequest{target: target, output: output}
}

func (c *curlRequest) HTTPEvent(event proxy.HTTPEvent) {
	switch event.EventType {
	case proxy.RequestEventType:
		c.writeCurlRequest(event)
	case proxy.ResponseEventType:
		return
	default:
		return
	}
}

func (c *curlRequest) writeCurlRequest(event proxy.HTTPEvent) {
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
