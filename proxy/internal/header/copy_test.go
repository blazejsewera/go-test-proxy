package header_test

import (
	"net/http"
	"testing"

	"github.com/blazejsewera/go-test-proxy/proxy/internal/header"
	"github.com/blazejsewera/go-test-proxy/test/assert"
)

func TestCopyHeaders(t *testing.T) {
	t.Run("copies header as-is", func(t *testing.T) {
		src := http.Header{
			"Connection": []string{"keep-alive"},
			"Accept":     []string{"application/json", "text/plain"},
		}
		expected := http.Header{
			"Connection": []string{"keep-alive"},
			"Accept":     []string{"application/json", "text/plain"},
		}

		dst := http.Header{}
		header.Copy(dst, src)

		assert.DeepEqual(t, expected, dst)
	})
}
