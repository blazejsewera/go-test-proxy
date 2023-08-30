package assert

import (
	"fmt"
	"slices"
)

func HeaderContainsExpected(expected, actual map[string][]string) error {
	for key, expectedValues := range expected {
		actualValues, ok := actual[key]
		if !ok {
			return fmt.Errorf("header: key '%s' not found in actual", key)
		}
		for _, expectedValue := range expectedValues {
			ok = slices.Contains(actualValues, expectedValue)
			if !ok {
				return fmt.Errorf("header: expected value '%s' for key '%s' not found in actual", expectedValue, key)
			}
		}
	}
	return nil
}
