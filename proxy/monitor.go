package proxy

import (
	"encoding/json"
	"log"
)

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

	Status int `json:"status"`
}

type eventMonitor interface {
	HTTPEvent(event HTTPEvent)
}

type errorMonitor interface {
	Err(err error)
}

type Monitor interface {
	eventMonitor
	errorMonitor
}

type DefaultMonitor struct{}

var _ Monitor = DefaultMonitor{}

func (n DefaultMonitor) HTTPEvent(event HTTPEvent) {
	eventJSONBytes, err := json.MarshalIndent(event, "", "\t")
	if err != nil {
		n.Err(err)
		return
	}
	log.Println(string(eventJSONBytes))
}

func (n DefaultMonitor) Err(err error) {
	log.Println(err)
}
