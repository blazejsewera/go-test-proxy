package proxytest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

type TestServerBuilder struct {
	target string
}

func (b *TestServerBuilder) WithTarget(url string) *TestServerBuilder {
	b.target = url
	return b
}

func (b *TestServerBuilder) Build() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(b.target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "url parse: %s\n", err)
			return
		}
		clonedRequest := r.Clone(context.Background())
		clonedRequest.RequestURI = ""
		clonedRequest.URL.Host = targetURL.Host

		response, err := http.DefaultClient.Do(clonedRequest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "response: %s\n", err)
			return
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "response reading: %s\n", err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "response writing: %s\n", err)
			return
		}
	})

	return httptest.NewUnstartedServer(handler)
}

func Builder() *TestServerBuilder {
	return &TestServerBuilder{}
}
