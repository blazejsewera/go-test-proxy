package proxytest

import (
	"github.com/blazejsewera/go-test-proxy/urls"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

type TestServerBuilder struct {
	target string
}

func (b *TestServerBuilder) WithTarget(url string) *TestServerBuilder {
	log.Printf("b.target = %s\n", b.target)
	b.target = url
	return b
}

func (b *TestServerBuilder) Build() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL := urls.ForwardedURL(b.target, r.URL)
		log.Printf("b.target = %s; target = %s; r = %v", b.target, targetURL, r)

		r.Host = targetURL.Host
		r.RequestURI = ""
		r.URL = targetURL

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Printf("%v\n", targetURL)
			log.Printf("response: %s\n", err)
			return
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("response reading: %s\n", err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			log.Printf("response writing: %s\n", err)
			return
		}
	})

	return httptest.NewUnstartedServer(handler)
}

func Builder() *TestServerBuilder {
	return &TestServerBuilder{}
}
