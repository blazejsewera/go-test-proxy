package requests

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

const assertionPrefix = "assert requests equal: "

func AssertEqualExcludingHost(t testing.TB, expected, actual *http.Request) {
	t.Helper()

	assertURLEqualExcludingHost(t, expected.URL, actual.URL)
	assertHeadersEqual(t, expected, actual)
	assertBodyEqual(t, expected, actual)
}

func assertURLEqualExcludingHost(t testing.TB, expected, actual *url.URL) {
	t.Helper()

	ok := expected.RawPath == actual.RawPath
	if !ok {
		assertionErrorf(t, "URL path not equal: %s, %s", expected.RawPath, actual.RawPath)
	}

	ok = expected.RawQuery == actual.RawQuery
	if !ok {
		assertionErrorf(t, "URL query not equal: %s, %s", expected.RawQuery, actual.RawQuery)
	}

	ok = expected.Scheme == actual.Scheme
	if !ok {
		assertionErrorf(t, "URL scheme not equal: %s, %s", expected.Scheme, actual.Scheme)
	}
}

func assertHeadersEqual(t testing.TB, expected, actual *http.Request) {
	t.Helper()
	ok := reflect.DeepEqual(expected.Header, actual.Header)
	if !ok {
		assertionErrorf(t, "headers not equal: %v, %v", expected.Header, actual.Header)
	}
}

func assertBodyEqual(t testing.TB, expected, actual *http.Request) {
	t.Helper()
	expectedBodyBytes, err := io.ReadAll(expected.Body)
	if err != nil {
		assertionFatalf(t, "read expected body: %s", err)
	}
	actualBodyBytes, err := io.ReadAll(actual.Body)
	if err != nil {
		assertionFatalf(t, "read actual body: %s", err)
	}
	expectedBody := string(expectedBodyBytes)
	actualBody := string(actualBodyBytes)

	ok := expectedBody == actualBody
	if !ok {
		assertionErrorf(t, "body not equal: %s, %s", expectedBody, actualBody)
	}
}

func assertionErrorf(t testing.TB, format string, args ...any) {
	t.Helper()
	t.Errorf(assertionPrefix+format, args...)
}

func assertionFatalf(t testing.TB, format string, args ...any) {
	t.Helper()
	t.Errorf(assertionPrefix+format, args...)
}
