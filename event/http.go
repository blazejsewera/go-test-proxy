package event

type Type string

const (
	RequestEventType  Type = "request"
	ResponseEventType Type = "response"
)

type HTTP struct {
	EventType Type `json:"eventType"`

	// common data

	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`

	// request data

	Method string `json:"method"`
	Path   string `json:"path"`
	Query  string `json:"query"`

	// response data

	Status int `json:"status"`
}
