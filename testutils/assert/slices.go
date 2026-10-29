package assert

import "testing"

func Empty[T any](t testing.TB, s []T) {
	t.Helper()
	if len(s) != 0 {
		t.Errorf("assert: not empty: %v\n", s)
	}
}
