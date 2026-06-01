package monitor

import (
	"fmt"
	"io"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/event"
)

type curlRequest struct {
	target string
	output io.Writer
	cfmt   *colorfmt.Fmt
}

func NewCurlRequestMonitor(target string, cfmt *colorfmt.Fmt) Monitor {
	return &curlRequest{target: target, cfmt: cfmt}
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
	c.cfmt.Cprintf(colorfmt.Faint, colorfmt.Green, "curl -X %s%s%s %s\n",
		e.Method,
		headerToCurl(e.Header),
		bodyToCurl(e.Body),
		urlPartsToCurl(c.target, e.Path, e.Query))
}

func headerToCurl(header map[string][]string) string {
	result := ""
	for key, values := range header {
		for _, value := range values {
			formattedHeader := fmt.Sprintf("%s: %s", key, value)
            result += fmt.Sprintf(" -H %q", formattedHeader)
		}
	}
	return result
}

func bodyToCurl(body string) string {
	if body == "" {
		return ""
	}
	return fmt.Sprintf(" -d %q", body)
}

func urlPartsToCurl(host string, path string, query string) string {
	if query == "" {
		return host + path
	}
	return fmt.Sprintf("%s%s?%s", host, path, query)
}

func (c *curlRequest) Err(error) {}
