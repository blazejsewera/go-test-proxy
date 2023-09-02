package monitor

import (
	"github.com/blazejsewera/go-test-proxy/proxy"
	"io"
	"os"
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
}

func (c *CurlRequestVoidResponse) Err(err error) {
}
