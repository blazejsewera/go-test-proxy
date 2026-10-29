---
name: mocking
description: Implementing a new handler function that returns a hardcoded (mocked) response.
---

# Mocking

To quickly mock endpoints with `gotestproxy`,
implement and add a new handler function to the `proxy.Builder`
with a `WithMock` or `WithMockGroup` method.

Go to `main.go` and add a mock that intercepts traffic on that path
regardless of enabled mock groups:

```go
package main

// ...

func main() {
	proxy.NewBuilder().
		// ...
		WithMock("/mockedPath", customHandlerFunc).
		// ...
		Build()
}
```

Or better yet, add a mock group:

```go
package main

// ...

func main() {
	proxy.NewBuilder().
		// ...
		WithMockGroup("group1",
			proxy.Mock{
				RoutePattern: "/someRoute",
				HandlerFunc:  someHandlerFunc,
			},
			proxy.Mock{
				RoutePattern: "/someOtherRoute",
				HandlerFunc:  someOtherHandlerFunc,
			}).
		WithMockGroup("group2",
			proxy.Mock{
				RoutePattern: "/anotherRoute",
				HandlerFunc:  anotherHandlerFunc,
			}).
		// ...
		Build()
}
```

You should only add a new file with a mock handler, preferably in the `mock` package and edit `main.go`.
There is no need to touch any other file.

## Example

A recommended pattern is to define handlers in a new file in the `mock` package.
For example:

```go
// file: mock/endpoint.go
package mock

import (
	"net/http"
)

const (
	EndpointRoute = "/endpoint"
)

func EndpointHandlerNormal(w http.ResponseWriter, r *http.Request) {
	requestQueryParam := r.URL.Query().Get("param")
	body := []byte(`{"key": "static value", "param": "` + requestQueryParam + `"}`)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		panic(err)
	}
}

func EndpointHandlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
```

```go
// file: main.go
package main

import "github.com/blazejsewera/go-test-proxy/mock"

// ...

func main() {
	proxy.NewBuilder().
		// ...
		WithMockGroup("endpointNormal",
			proxy.Mock{
				RoutePattern: mock.EndpointRoute,
				HandlerFunc:  mock.EndpointHandlerNormal,
			}).
		WithMockGroup("endpointNotFound",
			proxy.Mock{
				RoutePattern: mock.EndpointRoute,
				HandlerFunc:  mock.EndpointHandlerNotFound,
			}).
		// ...
		Build()
}
```

Take a note of the mock groups added in the `main.go` file, and when you need to run the `gotestproxy`,
do it with one of the groups enabled (here: `-mockGroups=endpointNormal` or `-mockGroups=endpointNotFound`),
depending on your needs.
