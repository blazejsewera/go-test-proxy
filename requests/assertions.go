package requests

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

const assertionPrefix = "assert requests: "

func AssertEqualExcludingHost(expected, actual *http.Request) error {
	err1 := assertURLEqualExcludingHost(expected.URL, actual.URL)
	err2 := assertHeadersEqual(expected, actual)
	err3 := assertBodyEqual(expected, actual)
	return errors.Join(err1, err2, err3)
}

func assertURLEqualExcludingHost(expected, actual *url.URL) error {
	var err1, err2, err3 error

	if expected.RawPath != actual.RawPath {
		err1 = assertionErrorf("URL path not equal: %s, %s", expected.RawPath, actual.RawPath)
	}

	if expected.RawQuery != actual.RawQuery {
		err2 = assertionErrorf("URL query not equal: %s, %s", expected.RawQuery, actual.RawQuery)
	}

	if expected.Scheme != actual.Scheme {
		err3 = assertionErrorf("URL scheme not equal: %s, %s", expected.Scheme, actual.Scheme)
	}

	return errors.Join(err1, err2, err3)
}

func assertHeadersEqual(expected, actual *http.Request) error {
	ok := reflect.DeepEqual(expected.Header, actual.Header)
	if !ok {
		return assertionErrorf("headers not equal: %v, %v", expected.Header, actual.Header)
	}
	return nil
}

func assertBodyEqual(expected, actual *http.Request) error {
	var err1, err2, err3 error

	expectedBodyBytes, err := io.ReadAll(expected.Body)
	if err != nil {
		err1 = assertionFatalf("read expected body: %s", err)
	}

	actualBodyBytes, err := io.ReadAll(actual.Body)
	if err != nil {
		err2 = assertionFatalf("read actual body: %s", err)
	}

	expectedBody := string(expectedBodyBytes)
	actualBody := string(actualBodyBytes)
	if expectedBody != actualBody {
		err3 = assertionErrorf("body not equal: %s, %s", expectedBody, actualBody)
	}

	return errors.Join(err1, err2, err3)
}

func assertionErrorf(format string, args ...any) error {
	return fmt.Errorf(assertionPrefix+format, args...)
}

func assertionFatalf(format string, args ...any) error {
	return fmt.Errorf(assertionPrefix+"fatal: "+format, args...)
}
