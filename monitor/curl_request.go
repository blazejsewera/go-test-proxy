package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"io"
	"os"
	"strconv"
)

type CurlRequestVoidResponse struct {
	target string
	output io.Writer
}

func NewCurlRequestVoidResponseMonitor(target string, output io.Writer) *CurlRequestVoidResponse {
	return &CurlRequestVoidResponse{target: target, output: output}
}

func NewCurlRequestVoidResponseMonitorToStdOut(target string) *CurlRequestVoidResponse {
	return NewCurlRequestVoidResponseMonitor(target, os.Stdout)
}

var _ proxy.Monitor = (*CurlRequestVoidResponse)(nil)

func (c *CurlRequestVoidResponse) HTTPEvent(event proxy.HTTPEvent) {
	switch event.EventType {
	case proxy.RequestEventType:
		c.writeCurlRequest(event)
	case proxy.ResponseEventType:
		return
	default:
		return
	}
}

func (c *CurlRequestVoidResponse) Err(err error) {
	_, errW := fmt.Fprintf(os.Stderr, "[PROXY ERROR]: %s\n", err)
	if errW != nil {
		panic("cannot write to stderr")
	}
}

func (c *CurlRequestVoidResponse) writeCurlRequest(event proxy.HTTPEvent) {
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
