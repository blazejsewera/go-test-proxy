package monitor

import (
	"github.com/blazejsewera/go-test-proxy/event"
)

type Monitor interface {
	// Name is used for enabling or disabling the monitor based on -monitor=name1,name2 arg.
	// The monitor is always included in the CombinedMonitor when the Name is empty.
	Name() string

	// HTTPEvent is called whenever a request or a response reaches the Monitor.
	// Implement it so that it produces some useful information to the console.
	// Use the colorfmt package to produce colorful output.
	HTTPEvent(e event.HTTP)
}

type NopMonitor struct{}

var _ Monitor = NopMonitor{}

func (n NopMonitor) Name() string {
	return "nop"
}

func (n NopMonitor) HTTPEvent(event.HTTP) {}
