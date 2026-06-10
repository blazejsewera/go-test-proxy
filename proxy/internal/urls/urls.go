package urls

import (
	"net/url"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
)

func ForwardedURL(rawTargetHost string, incomingURL *url.URL) *url.URL {
	targetURL, err := url.Parse(rawTargetHost)
	if err != nil {
		log.Fatalf("target host: url parse: %v\n", err)
	}

	targetURL.Path = incomingURL.Path
	targetURL.RawPath = incomingURL.RawPath
	targetURL.RawQuery = incomingURL.RawQuery
	targetURL.Fragment = incomingURL.Fragment
	targetURL.RawFragment = incomingURL.RawFragment

	return targetURL
}
