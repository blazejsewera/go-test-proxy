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
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "\n===== REQUEST =====\n")
	query := ""
	if e.Query != "" {
		query = fmt.Sprintf("?%s", e.Query)
	}
	c.cfmt.Cprintf(colorfmt.Bold, colorfmt.BrightBlue, "%s ", e.Method)
	c.cfmt.Cprintf(colorfmt.Underline, colorfmt.BrightBlue, "%s%s\n", e.Path, query)
	c.cfmt.Cprintf(colorfmt.Italic, colorfmt.Blue, "Target-Host: %s\n", c.target)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Blue, "--HEADER--\n")
	c.printHeader(e.Header, colorfmt.Blue)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Blue, "--BODY--\n")
	c.printBody(e.Body, colorfmt.Blue)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "===================\n")
}

func (c *console) printResponse(e event.HTTP) {
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "\n===== RESPONSE =====\n")
	c.printStatus(e.Status)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.BrightWhite, "--HEADER--\n")
	c.printHeader(e.Header, colorfmt.BrightWhite)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.BrightWhite, "--BODY--\n")
	c.printBody(e.Body, colorfmt.BrightWhite)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "====================\n")
}

func (c *console) printStatus(status int) {
	c.cfmt.Cprintf(colorfmt.Bold, colorfmt.BrightWhite, "STATUS: ")

	color := colorfmt.BrightWhite
	statusClass := status / 100
	switch statusClass {
	case 1:
		color = colorfmt.BrightBlue
	case 2:
		color = colorfmt.BrightGreen
	case 3:
		color = colorfmt.BrightBlue
	case 4, 5:
		color = colorfmt.BrightRed
	}

	c.cfmt.Cprintf(colorfmt.Bold, color, "%d\n", status)
}

func (c *console) printBody(body string, color colorfmt.Color) {
	if body == "" {
		c.cfmt.Cprint(colorfmt.Normal, color, "<empty>\n")
	} else {
		c.cfmt.Cprint(colorfmt.Normal, color, body+"\n")
	}
}

func (c *console) printHeader(h map[string][]string, color colorfmt.Color) {
	for key, values := range h {
		for _, value := range values {
			c.cfmt.Cprintf(colorfmt.Bold, color, "%s: ", key)
			c.cfmt.Cprintf(colorfmt.Normal, color, "%s\n", value)
		}
	}
}

func (c *console) Err(error) {}
