package requests

import (
	"fmt"
	"net/http"
	"strings"
)

func MustMakeNewRequest(method, url, body string, headers map[string]string) *http.Request {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(fmt.Errorf("construct new request: %s %s: %s", method, url, err))
		return nil
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}
