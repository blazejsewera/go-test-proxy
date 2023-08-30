package proxy

type EventType string

const (
	RequestEventType  EventType = "request"
	ResponseEventType EventType = "response"
)

type HTTPEvent struct {
	EventType EventType `json:"eventType"`

	// common data

	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`

	// request data

	Method string `json:"method"`
	Path   string `json:"path"`
	Query  string `json:"query"`

	// response data

	Status uint `json:"status"`
}

type Monitor interface {
	HTTPEvent(e HTTPEvent)
}

type NoopMonitor struct{}

var _ Monitor = NoopMonitor{}

func (n NoopMonitor) HTTPEvent(HTTPEvent) {}
