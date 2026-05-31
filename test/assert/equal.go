package assert

import (
	"strconv"
	"testing"
)

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("assert: not equal: expected = %v; actual = %v\n", expected, actual)
	}
}

// EqualQuoted is the same as Equal,
// but quotes the expected and actual strings in the assertion error message
func EqualQuoted(t testing.TB, expected, actual string) {
	t.Helper()
	Equal(t, strconv.Quote(expected), strconv.Quote(actual))
}
