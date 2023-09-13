package monitor

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
)

type console struct {
	target string
}

var _ proxy.Monitor = (*console)(nil)

func NewConsoleMonitor(target string) proxy.Monitor {
	return &console{
		target: target,
	}
}

func (c *console) HTTPEvent(event proxy.HTTPEvent) {
	if event.EventType == proxy.RequestEventType {
		c.printRequest(event)
	} else if event.EventType == proxy.ResponseEventType {
		c.printResponse(event)
	}
}

func (c *console) printRequest(event proxy.HTTPEvent) {
	fmt.Print("\n===== REQUEST =====\n")
	query := ""
	if event.Query != "" {
		query = fmt.Sprintf("?%s", event.Query)
	}
	fmt.Printf("%s %s%s\n", event.Method, event.Path, query)
	fmt.Printf("Target-Host: %s\n", c.target)
	fmt.Print("--HEADER--\n")
	printHeader(event.Header)
	fmt.Print("--BODY--\n")
	printBody(event.Body)
	fmt.Print("===================\n")
}

func (c *console) printResponse(event proxy.HTTPEvent) {
	fmt.Print("\n===== RESPONSE =====\n")
	fmt.Printf("STATUS: %d\n", event.Status)
	fmt.Print("--HEADER--\n")
	printHeader(event.Header)
	fmt.Print("--BODY--\n")
	printBody(event.Body)
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
