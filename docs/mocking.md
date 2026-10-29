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

You can safely add and enable mocks with the same `RoutePattern`, they will be handled gracefully with a warning.

## Example

A recommended pattern is to define handlers in a new file in the `mock` package.
You can look at the `mock/example.go` and `main.go` files or use the examples below:

```go
// file: mock/example.go
package mock

import (
	"encoding/json"
	"net/http"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
)

const (
	ExampleRoute = "/_example"
)

type response struct {
	ResponseValue string `json:"responseValue"`
}

func ExampleHandler200(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParam := r.URL.Query().Get("queryParam")
	res := response{queryParam}

	responseJSON, _ := json.Marshal(res)
	_, err := w.Write(responseJSON)
	if err != nil {
		log.Fatalf("error in ExampleHandler: %v\n", err)
	}
}

func ExampleHandler404(w http.ResponseWriter, _ *http.Request) {
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
		WithMockGroup("example200", proxy.Mock{RoutePattern: mock.ExampleRoute, HandlerFunc: mock.ExampleHandler200}).
		WithMockGroup("example404", proxy.Mock{RoutePattern: mock.ExampleRoute, HandlerFunc: mock.ExampleHandler200}).
		// ...
		Build()
}
```

Take a note of the mock groups added in the `main.go` file, and when you need to run the `gotestproxy`,
do it with one of the groups enabled (here: `-mockGroups=example200` or `-mockGroups=example404`),
depending on your needs.
