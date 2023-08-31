package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"os"
)

type Console struct {
	target string
}

var _ proxy.Monitor = (*Console)(nil)

func NewConsoleMonitor(target string) *Console {
	return &Console{target: target}
}

func (c *Console) HTTPEvent(event proxy.HTTPEvent) {
	if event.EventType == proxy.RequestEventType {
		c.printRequest(event)
	} else if event.EventType == proxy.ResponseEventType {
		c.printResponse(event)
	}
}

func (c *Console) printRequest(event proxy.HTTPEvent) {
	fmt.Print("\n\n==REQUEST==\n")
	query := ""
	if event.Query != "" {
		query = fmt.Sprintf("?%s", event.Query)
	}
	fmt.Printf("%s %s%s\n", event.Method, event.Path, query)
	fmt.Printf("Target-Host: %s\n", c.target)
	fmt.Print("--HEADER--\n")
	printHeader(event.Header)
	fmt.Print("--BODY--\n")
	fmt.Println(event.Body)
	fmt.Print("===========\n")
}

func (c *Console) printResponse(event proxy.HTTPEvent) {
	fmt.Print("\n\n==RESPONSE==\n")
	fmt.Printf("STATUS: %d\n", event.Status)
	fmt.Print("--HEADER--\n")
	printHeader(event.Header)
	fmt.Print("--BODY--\n")
	fmt.Println(event.Body)
	fmt.Print("===========\n")
}

func printHeader(h map[string][]string) {
	for key, values := range h {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}

func (c *Console) Err(err error) {
	_, errW := fmt.Fprintf(os.Stderr, "[PROXY ERROR]: %s\n", err)
	if errW != nil {
		panic("cannot write to stderr")
	}
}
