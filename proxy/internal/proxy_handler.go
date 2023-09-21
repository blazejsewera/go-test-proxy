package internal

import (
	"fmt"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/header"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/urls"
	"io"
	"net/http"
)

func ProxyHandler(monitor monitor.Monitor, url string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL := urls.ForwardedURL(url, r.URL)

		r.RequestURI = ""
		r.Host = targetURL.Host
		r.URL = targetURL

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			monitor.Err(fmt.Errorf("client request to target: %s", err))
			return
		}
		w.WriteHeader(response.StatusCode)
		_, err = io.Copy(w, response.Body)
		if err != nil {
			monitor.Err(fmt.Errorf("write response: %s", err))
			return
		}
		header.Copy(w.Header(), response.Header)
	}
}
