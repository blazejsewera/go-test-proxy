package monitor

import (
	"fmt"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/event"
)

type response struct {
	target string
	cfmt   *colorfmt.Fmt
}

var _ Monitor = (*response)(nil)

type request struct {
	target string
	cfmt   *colorfmt.Fmt
}

func (r *request) Name() string {
	return "request"
}

func (r *request) HTTPEvent(e event.HTTP) {
	if e.EventType == event.RequestEventType {
		r.printRequest(e)
	}
}

func (r *request) Err(error) {}

func (r *request) printRequest(e event.HTTP) {
	r.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "\n===== REQUEST =====\n")
	query := ""
	if e.Query != "" {
		query = fmt.Sprintf("?%s", e.Query)
	}
	r.cfmt.Cprintf(colorfmt.Bold, colorfmt.BrightBlue, "%s ", e.Method)
	r.cfmt.Cprintf(colorfmt.Underline, colorfmt.BrightBlue, "%s%s\n", e.Path, query)
	r.cfmt.Cprintf(colorfmt.Italic, colorfmt.Blue, "Target-Host: %s\n", r.target)
	r.cfmt.Cprint(colorfmt.Faint, colorfmt.Blue, "--HEADER--\n")
	printHeader(r.cfmt, e.Header, colorfmt.Blue)
	r.cfmt.Cprint(colorfmt.Faint, colorfmt.Blue, "--BODY--\n")
	printBody(r.cfmt, e.Body, colorfmt.Blue)
	r.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "===================\n")
}

func NewRequestMonitor(target string, cfmt *colorfmt.Fmt) Monitor {
	return &request{target, cfmt}
}

func NewResponseMonitor(target string, cfmt *colorfmt.Fmt) Monitor {
	return &response{
		target: target,
		cfmt:   cfmt,
	}
}

func (c *response) Name() string {
	return "response"
}

func (c *response) HTTPEvent(e event.HTTP) {
	if e.EventType == event.ResponseEventType {
		c.printResponse(e)
	}
}

func (c *response) printResponse(e event.HTTP) {
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "\n===== RESPONSE =====\n")
	printStatus(c.cfmt, e.Status)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.BrightWhite, "--HEADER--\n")
	printHeader(c.cfmt, e.Header, colorfmt.BrightWhite)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.BrightWhite, "--BODY--\n")
	printBody(c.cfmt, e.Body, colorfmt.BrightWhite)
	c.cfmt.Cprint(colorfmt.Faint, colorfmt.Base, "====================\n")
}

func printStatus(cfmt *colorfmt.Fmt, status int) {
	cfmt.Cprintf(colorfmt.Bold, colorfmt.BrightWhite, "STATUS: ")

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

	cfmt.Cprintf(colorfmt.Bold, color, "%d\n", status)
}

func printBody(cfmt *colorfmt.Fmt, body string, color colorfmt.Color) {
	if body == "" {
		cfmt.Cprint(colorfmt.Normal, color, "<empty>\n")
	} else {
		cfmt.Cprint(colorfmt.Normal, color, body+"\n")
	}
}

func printHeader(cfmt *colorfmt.Fmt, h map[string][]string, color colorfmt.Color) {
	for key, values := range h {
		for _, value := range values {
			cfmt.Cprintf(colorfmt.Bold, color, "%s: ", key)
			cfmt.Cprintf(colorfmt.Normal, color, "%s\n", value)
		}
	}
}

func (c *response) Err(error) {}
