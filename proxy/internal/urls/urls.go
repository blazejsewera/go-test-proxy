package urls

import (
	"net/url"
)

func ForwardedURL(targetURL *url.URL, incomingURL *url.URL) *url.URL {
	result, _ := url.Parse(targetURL.String()) // targetURL.Clone() is available from Go 1.27
	var path, rawPath string
	if targetURL.Path != "" && targetURL.Path != "/" {
		prefixedPath := targetURL.JoinPath(incomingURL.Path)
		path = prefixedPath.Path
		rawPath = prefixedPath.RawPath
	} else {
		path = incomingURL.Path
		rawPath = incomingURL.RawPath
	}

	result.Path = path
	result.RawPath = rawPath

	result.RawQuery = incomingURL.RawQuery
	result.Fragment = incomingURL.Fragment
	result.RawFragment = incomingURL.RawFragment

	return result
}
