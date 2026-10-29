package assert

import (
	"strings"
	"testing"
)

func StringContains(t testing.TB, expectedSubstring string, actual string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("assert: does not contain: expected = %v; actual = %v\n", expectedSubstring, actual)
	}
}
