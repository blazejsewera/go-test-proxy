package monitor

import (
	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/event"
)

type stdErr struct {
	cfmt *colorfmt.Fmt
}

func NewStdErrMonitor(cfmt *colorfmt.Fmt) Monitor {
	return &stdErr{cfmt}
}

func (e *stdErr) Err(err error) {
	e.cfmt.Cerrprintf(colorfmt.Normal, colorfmt.Red, "[PROXY ERROR]: %s\n", err)
}

func (e *stdErr) HTTPEvent(event.HTTP) {}
