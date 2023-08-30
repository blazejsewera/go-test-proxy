package request

import (
	"fmt"
	"net/http"
	"strings"
)

func Reference() *http.Request {
	return New("", "")
}

func New(baseURL, path string) *http.Request {
	headers := map[string]string{"X-Test-Header": "Test-Value"}
	return mustBuildRequest("GET", baseURL+path, "body", headers)
}

func mustBuildRequest(method, url, body string, headers map[string]string) *http.Request {
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
