package monitor

import (
	"fmt"
	"os"

	"github.com/blazejsewera/go-test-proxy/event"
)

type stdErr struct{}

func NewStdErrMonitor() Monitor {
	return &stdErr{}
}

func (e *stdErr) Err(err error) {
	_, errW := fmt.Fprintf(os.Stderr, "[PROXY ERROR]: %s\n", err)
	if errW != nil {
		panic("cannot write to stderr")
	}
}

func (e *stdErr) HTTPEvent(event.HTTP) {}
