package requests

import (
	"net/http"
	"strings"
	"testing"
)

func MustMakeNewRequest(t testing.TB, method, url, body string, headers map[string]string) *http.Request {
	t.Helper()
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("construct new request: %s %s: %s", method, url, err)
		return nil
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}
