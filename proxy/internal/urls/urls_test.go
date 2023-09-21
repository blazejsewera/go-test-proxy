package urls_test

import (
	"github.com/blazejsewera/go-test-proxy/proxy/internal/urls"
	"github.com/blazejsewera/go-test-proxy/test/assert"
	"net/url"
	"testing"
)

func TestForwardedURL(t *testing.T) {
	t.Run("target url has path, query, and fragment copied from incoming url "+
		"and scheme, user, host, and port copied from target host", func(t *testing.T) {
		incomingURL := URLMustParse(t, "http://user@host:1337/path?query=value#fragment")
		targetHost := "https://user1@targetHost:9001"
		expected := "https://user1@targetHost:9001/path?query=value#fragment"

		actual := urls.ForwardedURL(targetHost, incomingURL)

		assert.Equal(t, expected, actual.String())
	})
}

func URLMustParse(t testing.TB, rawURL string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("url parse: %s\n", err)
	}
	return parsed
}
