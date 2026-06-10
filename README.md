# Go Test Proxy

A Go proxy to plug between the frontend and the backend
or between backend services for last-resort testing.

## Why?

I only want to see the requests and responses,
without the necessity to open a debugger.

## Quick start

### Build

```sh
make build
```

or if you don't have Make:

```sh
go build github.com/blazejsewera/go-test-proxy/cmd/gotestproxy
```

### Run

```sh
./gotestproxy -target=https://example.com -port=8000 -color
```

### Point your client to the proxy

Depending on your project, you want to have something like this:

```yaml
backendUrl: "http://localhost:8000"
```

## Mocking certain endpoints

It is easy to quickly mock endpoints with Go Test Proxy,
simply add a new handler function to the proxy builder.

Go to [main.go](main.go) and add a mock that intercepts traffic on that path
regardless of enabled mock groups:

```go
package main

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

You can extract the `proxy.Mock` creation to a factory function being in another package, like `mock`.

Then rebuild the project and run it with `-allMocks` arg to enable all mock groups,
or `-mockGroups=group1,group2` arg to enable specified mock groups.

## Swapping the monitor implementations

Look at `builder.WithMonitor` and `monitor.Combine` functions.
The former lets you set any monitor adhering to the `proxy.Monitor` interface.
The latter lets you combine multiple monitors — they will be called one-by-one.

## Install

To install `gotestproxy` binary, simply run `make install`,
setting a `PREFIX` that is in your `PATH`.

```sh
PREFIX=<target_directory> make install
```

Uninstalling is also very simple:

```sh
PREFIX=<previously_set_prefix> make uninstall
```

Alternatively, you can simply remove the `gotestproxy` binary — it's self-contained,
so it's the only file to remove.

## Build for different targets

You can also build `gotestproxy` for different targets, like windows/amd64:

```sh
make build-windows-amd64
```
