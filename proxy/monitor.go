package proxy

type EventType string

const (
	Request  EventType = "request"
	Response EventType = "response"
)

type HttpEvent struct {
	EventType         EventType `json:"eventType"`
	CustomHandlerUsed bool      `json:"customHandlerUsed"`

	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`

	Method string `json:"method"`
	Path   string `json:"path"`
	Query  string `json:"query"`

	Status uint `json:"status"`
}

type Monitor interface {
	HttpEvent(e HttpEvent)
}

type NoopMonitor struct{}

var _ Monitor = NoopMonitor{}

func (n NoopMonitor) HttpEvent(HttpEvent) {}
