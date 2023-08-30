package assert

import (
	"encoding/json"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"testing"
)

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("assert: not equal: expected = %v; actual = %v\n", expected, actual)
	}
}

func JSONEqual(t testing.TB, expected, actual any) {
	t.Helper()
	expectedJSON := marshalIndent(expected)
	actualJSON := marshalIndent(actual)
	if expectedJSON != actualJSON {
		t.Errorf("assert: not equal:\nexpected = %v\nactual = %v\n", expectedJSON, actualJSON)
	}
}

func marshalIndent(v any) string {
	return string(must.Succeed(json.MarshalIndent(v, "", "\t")))
}
