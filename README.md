# Go Test Proxy

A Go proxy to plug between the frontend and the backend
or between backend services for mocking endpoints or last-resort testing.

I find it useful for frontend development, when the API specs are agreed upon,
but there is no backend implementation yet.

This proxy also proved very helpful when developing mobile apps,
especially when you need to make changes in a project with no robust integration tests.

In mobile app development and between backends,
there is usually no convenient network traffic monitoring like in the browser devtools.
You can launch as many proxies as you want between any two components
that communicate with HTTP.

## Prerequisites

You'll need Go installed:

- [Official installation docs](https://go.dev/doc/install)
- `sudo apt install golang-go` for Ubuntu
- `brew install go` for macOS

Make sure that your Go version is at least 1.24 (run `go version`).

## Quick start

Absolute one-liner:

```sh
go run github.com/blazejsewera/go-test-proxy -target=https://example.com -port=8000 -color
```

### Build

```sh
make build
```

or if you don't have Make:

```sh
go build -o gotestproxy github.com/blazejsewera/go-test-proxy
```

### Run

```sh
./gotestproxy -target=https://example.com -port=8000 -color
```

To print help:

```sh
./gotestproxy -h
```

### Point your frontend client to the proxy

You'll need to change the URL in you frontend configuration
to point to the proxy instead of a real backend instance.

If you run your frontend on the host network (without Docker):

```diff
- backendUrl: "https://real-backend-url.in-the.cloud"
+ backendUrl: "http://localhost:8000"
```

If you run your frontend in a Docker network, but the proxy is in the host network,
the easiest way is to run `ip address`, `ifconfig`, or `ipconfig`,
and get the network card address, e.g.:

```
$ ifconfig
en0: flags=[some flags]
        options=[some options]
        ether xx:xx:xx:xx:xx:xx
        inet6 xxxx::xxxx:xxxx:xxxx:xxxx%en0 prefixlen 64 secured scopeid 0xe 
        inet 192.168.1.101 netmask 0xffffff00 broadcast 192.168.1.255
        media: autoselect
        status: active
```

because `inet` is `192.168.1.101`, we'll set the backend URL as such:

```diff
- backendUrl: "https://real-backend-url.in-the.cloud"
+ backendUrl: "http://192.168.1.101:8000"
```

and run `./gotestproxy -target=https://real-backend-url.in-the.cloud -port 8000 -color`.

## Mocking certain endpoints

It is easy to quickly mock endpoints with Go Test Proxy,
simply add a new handler function to the proxy builder.

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

Then rebuild the project and run it with `-allMocks` arg to enable all mock groups,
or `-mockGroups=group1,group2` arg to enable specified mock groups.

A pattern I find useful is to define handlers in a new file in the `mock` package.
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

and run the proxy with `-mockGroups=endpointNormal` or `-mockGroups=endpointNotFound` depending on your needs.

## Using monitors

The monitors let you print information to the console whenever a request or response goes through the proxy.

### Available monitors

- request / response monitor: shows basic HTTP info: method, headers, body, status
- curl request monitor: shows the request as a curl command that you can copy and paste into another program
  like Bruno or Insomnia, or run directly in a shell

### Adjusting which monitors are enabled

Look at `builder.WithMonitor` and `monitor.Combine` in `main.go`.
The former lets you set any monitor adhering to the `monitor.Monitor` interface.
The latter lets you combine multiple monitors — they will be called sequentially
on every intercepted request and response.

### Writing new monitors

Implement the `monitor.Monitor` interface,
then wire the new monitor up in `builder.WithMonitor(monitor.Combine([oldMonitors...], newMonitor))`.
Look at the `monitor.console` implementation for inspiration.
I encourage you to use the `colorfmt` package when writing to the console,
so that you'll have a pretty, colorful output.

## Install

To install `gotestproxy` binary, simply run `make install`,
setting the `PREFIX` env variable to a directory that is in your `PATH` env variable.

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

## Why Go

I used Go because it allowed me to avoid any external dependencies
and rely only on the standard library.
This in turn made this project more plug-and-play,
easier to set up, more portable, and very resistant to supply chain attacks.
