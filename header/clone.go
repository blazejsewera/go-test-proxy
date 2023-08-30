package header

import (
	"net/http"
)

func CloneToResponseWriter(header http.Header, w http.ResponseWriter) {
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

func Clone(source map[string][]string) map[string][]string {
	target := make(map[string][]string)
	for key, sourceValues := range source {
		targetValues := make([]string, len(sourceValues))
		copy(targetValues, sourceValues)
		target[key] = targetValues
	}
	return target
}
