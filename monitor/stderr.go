package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"os"
)

type stdErr struct{}

func NewStdErrMonitor() proxy.Monitor {
	return &stdErr{}
}

func (e *stdErr) Err(err error) {
	_, errW := fmt.Fprintf(os.Stderr, "[PROXY ERROR]: %s\n", err)
	if errW != nil {
		panic("cannot write to stderr")
	}
}

func (e *stdErr) HTTPEvent(proxy.HTTPEvent) {}
