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
./gotestproxy --target=https://example.com --port=8000
```

### Point your client to the proxy

Depending on your project, you want to have something like this:

```yaml
backendUrl: "http://localhost:8000"
```

## Mocking certain endpoints

It is easy to quickly mock endpoints with Go Test Proxy,
simply add a new handler function to the proxy builder.

Go to [main.go](cmd/gotestproxy/main.go) and invoke `builder.WithHandlerFunc("/mockedPath", customFunc)`.
Then rebuild the project and run it.

## Swapping the monitor implementations

Look at `builder.WithMonitor` and `monitor.Combine` functions.
The former lets you set any monitor adhering to the `proxy.Monitor` interface.
The latter lets you combine multiple monitors — they will be called one-by-one.
