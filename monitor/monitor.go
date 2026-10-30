package monitor

import (
	"github.com/blazejsewera/go-test-proxy/event"
)

type Monitor interface {
	// Name is used for enabling or disabling the monitor based on -monitor=name1,name2 arg.
	Name() string

	// HTTPEvent is called whenever a request or a response reaches the Monitor.
	// The implementation should produce useful output to stdout.
	// For writing to stdout use *colorfmt.Fmt.
	HTTPEvent(e event.HTTP)
}

type NopMonitor struct{}

var _ Monitor = NopMonitor{}

func (n NopMonitor) Name() string {
	return "nop"
}

func (n NopMonitor) HTTPEvent(event.HTTP) {}
