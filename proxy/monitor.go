package proxy

type EventType string

const (
	RequestEventType  EventType = "request"
	ResponseEventType EventType = "response"
)

type HttpEvent struct {
	EventType EventType `json:"eventType"`

	// common data

	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`

	// request data

	Method string `json:"method"`
	Path   string `json:"path"`
	Query  string `json:"query"`

	// response data

	Status            uint `json:"status"`
	CustomHandlerUsed bool `json:"customHandlerUsed"`
}

type Monitor interface {
	HttpEvent(e HttpEvent)
}

type NoopMonitor struct{}

var _ Monitor = NoopMonitor{}

func (n NoopMonitor) HttpEvent(HttpEvent) {}
