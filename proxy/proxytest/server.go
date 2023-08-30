package proxytest

import "net/http/httptest"

type TestServer struct {
	*httptest.Server
}
