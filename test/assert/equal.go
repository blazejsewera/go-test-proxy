package assert

import (
	"reflect"
	"testing"
)

const assertEqualFailedFormat = "assert: not equal: expected = %v; actual = %v\n"

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf(assertEqualFailedFormat, expected, actual)
	}
}

func DeepEqual(t testing.TB, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(assertEqualFailedFormat, expected, actual)
	}
}
