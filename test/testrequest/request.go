package testrequest

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	MethodGet     = "GET"
	ReferenceBody = "body"
)

func New(baseURL, path string) *http.Request {
	return mustBuildRequest(MethodGet, baseURL+path, ReferenceBody, ReferenceHeader())
}

func mustBuildRequest(method, url, body string, header http.Header) *http.Request {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(fmt.Errorf("construct new request: %s %s: %s", method, url, err))
		return nil
	}
	for key, values := range header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return request
}

func ReferenceRequest() *http.Request {
	return New("", "")
}

func ReferenceHeader() http.Header {
	return http.Header{"X-Test-Header": []string{"Test-Value"}}
}

func ReferenceResponseHeader() http.Header {
	return http.Header{"X-Response-Test-Header": []string{"Test-Value"}}
}
