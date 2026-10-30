package event

type Type string

const (
	RequestEventType  Type = "request"
	ResponseEventType Type = "response"
)

type HTTP struct {
	// EventType is either RequestEventType or ResponseEventType
	EventType Type

	// Header is present when EventType is either RequestEventType or ResponseEventType
	Header map[string][]string
	// Body is present when EventType is either RequestEventType or ResponseEventType
	Body string

	// Method is present only when EventType is RequestEventType
	Method string
	// Path is present only when EventType is RequestEventType
	Path string
	// Query is present only when EventType is RequestEventType
	Query string

	// Status is present only when EventType is ResponseEventType
	Status int
}
