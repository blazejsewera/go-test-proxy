package assert

import "testing"

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("assert: not equal: expected = %v; actual = %v\n", expected, actual)
	}
}
