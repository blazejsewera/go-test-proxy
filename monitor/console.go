package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/event"
)

type console struct {
	target string
}

var _ Monitor = (*console)(nil)

func NewConsoleMonitor(target string) Monitor {
	return &console{
		target: target,
	}
}

func (c *console) HTTPEvent(e event.HTTP) {
	if e.EventType == event.RequestEventType {
		c.printRequest(e)
	} else if e.EventType == event.ResponseEventType {
		c.printResponse(e)
	}
}

func (c *console) printRequest(e event.HTTP) {
	fmt.Print("\n===== REQUEST =====\n")
	query := ""
	if e.Query != "" {
		query = fmt.Sprintf("?%s", e.Query)
	}
	fmt.Printf("%s %s%s\n", e.Method, e.Path, query)
	fmt.Printf("Target-Host: %s\n", c.target)
	fmt.Print("--HEADER--\n")
	printHeader(e.Header)
	fmt.Print("--BODY--\n")
	printBody(e.Body)
	fmt.Print("===================\n")
}

func (c *console) printResponse(e event.HTTP) {
	fmt.Print("\n===== RESPONSE =====\n")
	fmt.Printf("STATUS: %d\n", e.Status)
	fmt.Print("--HEADER--\n")
	printHeader(e.Header)
	fmt.Print("--BODY--\n")
	printBody(e.Body)
	fmt.Print("====================\n")
}

func printBody(body string) {
	if body == "" {
		fmt.Println("<empty>")
	} else {
		fmt.Println(body)
	}
}

func printHeader(h map[string][]string) {
	for key, values := range h {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}

func (c *console) Err(error) {}
