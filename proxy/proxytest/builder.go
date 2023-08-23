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
			fmt.Fprintln(os.Stderr, "url parse error")
			return
		}
		clonedRequest := r.Clone(context.Background())
		clonedRequest.URL.Host = targetURL.Host
		fmt.Printf("%v", clonedRequest)

		response, err := http.DefaultClient.Do(clonedRequest)
		if err != nil {
			fmt.Fprintln(os.Stderr, "response error")
			return
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "response reading error")
			return
		}
		_, _ = w.Write(bytes)
	})

	return httptest.NewUnstartedServer(handler)
}

func Builder() *TestServerBuilder {
	return &TestServerBuilder{}
}
