package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/header"
	"github.com/blazejsewera/go-test-proxy/proxy/internal/urls"
)

func ProxyHandler(monitor monitor.Monitor, targetURL *url.URL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		forwardedURL := urls.ForwardedURL(targetURL, r.URL)

		r.RequestURI = ""
		r.Host = forwardedURL.Host
		r.URL = forwardedURL

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
