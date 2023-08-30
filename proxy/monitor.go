package proxy

type EventType string

const (
	RequestEventType  EventType = "request"
	ResponseEventType EventType = "response"
)

type HttpEvent struct {
	EventType         EventType `json:"eventType"`
	CustomHandlerUsed bool      `json:"customHandlerUsed"`

	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`

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
