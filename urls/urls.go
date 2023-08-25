package urls

import (
	"log"
	"net/url"
)

func ForwardedURL(rawTargetHost string, incomingURL *url.URL) *url.URL {
	targetURL, err := url.Parse(rawTargetHost)
	if err != nil {
		log.Fatalln("target host: url parse:", err)
	}

	targetURL.Path = incomingURL.Path
	targetURL.RawPath = incomingURL.RawPath
	targetURL.RawQuery = incomingURL.RawQuery
	targetURL.Fragment = incomingURL.Fragment
	targetURL.RawFragment = incomingURL.RawFragment

	return targetURL
}
