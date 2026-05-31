package monitor

import (
	"fmt"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/event"
)

type console struct {
	target string
	cfmt   *colorfmt.Fmt
}

var _ Monitor = (*console)(nil)

func NewConsoleMonitor(target string, cfmt *colorfmt.Fmt) Monitor {
	return &console{
		target: target,
		cfmt:   cfmt,
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
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "\n===== REQUEST =====\n")
	query := ""
	if e.Query != "" {
		query = fmt.Sprintf("?%s", e.Query)
	}
	c.cfmt.Cprintf(colorfmt.Normal, colorfmt.Base, "%s %s%s\n", e.Method, e.Path, query)
	c.cfmt.Cprintf(colorfmt.Normal, colorfmt.Base, "Target-Host: %s\n", c.target)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "--HEADER--\n")
	c.printHeader(e.Header)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "--BODY--\n")
	c.printBody(e.Body)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "===================\n")
}

func (c *console) printResponse(e event.HTTP) {
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "\n===== RESPONSE =====\n")
	c.cfmt.Cprintf(colorfmt.Normal, colorfmt.Base, "STATUS: %d\n", e.Status)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "--HEADER--\n")
	c.printHeader(e.Header)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "--BODY--\n")
	c.printBody(e.Body)
	c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "====================\n")
}

func (c *console) printBody(body string) {
	if body == "" {
		c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, "<empty>\n")
	} else {
		c.cfmt.Cprint(colorfmt.Normal, colorfmt.Base, body+"\n")
	}
}

func (c *console) printHeader(h map[string][]string) {
	for key, values := range h {
		for _, value := range values {
			c.cfmt.Cprintf(colorfmt.Normal, colorfmt.Base, "%s: %s\n", key, value)
		}
	}
}

func (c *console) Err(error) {}
