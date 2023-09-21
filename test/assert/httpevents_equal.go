package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/monitor/event"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"testing"
)

func HTTPEventListEqual(t testing.TB, expected, actual []event.HTTP) {
	t.Helper()

	if len(expected) != len(actual) {
		expectedJSON := marshalIndent(expected)
		actualJSON := marshalIndent(actual)

		t.Errorf("assert: http event lists are of different lengths")
		t.Errorf("expected = %v\nactual = %v", expectedJSON, actualJSON)
		return
	}

	for i, expectedEvent := range expected {
		actualEvent := actual[i]
		HTTPEventsEqual(t, expectedEvent, actualEvent)
	}
}

func HTTPEventsEqual(t testing.TB, expected, actual event.HTTP) {
	t.Helper()
	errs := make([]error, 5)

	errs = append(errs, assertString("eventType", string(expected.EventType), string(actual.EventType)))
	errs = append(errs, assertString("body", expected.Body, actual.Body))
	errs = append(errs, assertString("method", expected.Method, actual.Method))
	errs = append(errs, assertString("path", expected.Path, actual.Path))
	errs = append(errs, assertString("query", expected.Query, actual.Query))
	errs = append(errs, assertInt("status", expected.Status, actual.Status))
	errs = append(errs, assertHeaderContainsExpected(expected.Header, actual.Header))

	if err := errors.Join(errs...); err != nil {
		expectedJSON := marshalIndent(expected)
		actualJSON := marshalIndent(actual)

		t.Errorf("http events assert: %s\n", err)
		t.Errorf("expected = %v\nactual = %v\n", expectedJSON, actualJSON)
	}
}

func assertHeaderContainsExpected(expected map[string][]string, actual map[string][]string) error {
	return HeaderContainsExpectedToErr(expected, actual)
}

func assertString(name, expected, actual string) error {
	if expected != actual {
		return fmt.Errorf("%s: '%s' not equal to '%s'", name, expected, actual)
	}
	return nil
}

func assertInt(name string, expected, actual int) error {
	if expected != actual {
		return fmt.Errorf("%s: %d not equal to %d", name, expected, actual)
	}
	return nil
}

func marshalIndent(v any) string {
	return string(must.Succeed(json.MarshalIndent(v, "", "\t")))
}
