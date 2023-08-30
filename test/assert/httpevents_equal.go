package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"github.com/blazejsewera/go-test-proxy/test/must"
	"testing"
)

func HTTPEventsEqual(t testing.TB, expected, actual []proxy.HTTPEvent) {
	t.Helper()
	expectedJSON := marshalIndent(expected)
	actualJSON := marshalIndent(actual)
	if len(expected) != len(actual) {
		t.Errorf("assert: http event lists are of different lengths")
		t.Errorf("expected = %v\nactual = %v", expectedJSON, actualJSON)
		return
	}
	errs := make([]error, 5)

	for i, expectedEvent := range expected {
		actualEvent := actual[i]
		errs = append(errs, assertString("eventType", string(expectedEvent.EventType), string(actualEvent.EventType)))
		errs = append(errs, assertString("body", expectedEvent.Body, actualEvent.Body))
		errs = append(errs, assertString("method", expectedEvent.Method, actualEvent.Method))
		errs = append(errs, assertString("path", expectedEvent.Path, actualEvent.Path))
		errs = append(errs, assertString("query", expectedEvent.Query, actualEvent.Query))
		errs = append(errs, assertUint("status", expectedEvent.Status, actualEvent.Status))
		errs = append(errs, assertHeaderContainsExpected("header", expectedEvent.Header, actualEvent.Header))
	}

	if err := errors.Join(errs...); err != nil {
		t.Errorf("http events assert: %s", err)
		t.Errorf("expected = %v\nactual = %v", expectedJSON, actualJSON)
	}
}

func assertHeaderContainsExpected(name string, expected map[string][]string, actual map[string][]string) error {
	for key, expectedValues := range expected {
		actualValues, ok := actual[key]
		if !ok {
			return fmt.Errorf("%s: key '%s' not found in actual headers", name, key)
		}
		expectedValuesJSON := marshalIndent(expectedValues)
		actualValuesJSON := marshalIndent(actualValues)

		if expectedValuesJSON != actualValuesJSON {
			return fmt.Errorf("%s: not equal for key '%s'\nexpected = %s\nactual = %s", name, key, expectedValuesJSON, actualValuesJSON)
		}
	}
	return nil
}

func assertString(name, expected, actual string) error {
	if expected != actual {
		return fmt.Errorf("%s: '%s' not equal to '%s'", name, expected, actual)
	}
	return nil
}

func assertUint(name string, expected, actual uint) error {
	if expected != actual {
		return fmt.Errorf("%s,%d not equal to %d", name, expected, actual)
	}
	return nil
}

func marshalIndent(v any) string {
	return string(must.Succeed(json.MarshalIndent(v, "", "\t")))
}
