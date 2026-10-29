# Go Test Proxy

A proxy to plug between the frontend and the backend,
or between backend services. Useful for mocking endpoints or last-resort testing.

<details>
<summary>Intended use case</summary>

I find it useful for frontend development, when the API specs are agreed upon,
but there is no backend implementation yet.

It also proved very helpful when developing mobile apps,
especially when you need to make changes in a project with no robust integration tests.

In mobile app development and between backends,
there is usually no convenient network traffic monitoring like in the browser devtools.
You can launch as many proxies as you want between any two components
that communicate with HTTP.
</details>

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

### Point your frontend client to a `gotestproxy` instance

You'll need to change the URL in you frontend configuration
to point to a `gotestproxy` instance instead of a real backend instance.
The `backendUrl` here is for illustration purposes;
you might have something different, but look for the place where you have your real backend URL configured.

If you run your frontend on the host network (without Docker):

```diff
-backendUrl: "https://real-backend-url.in-the.cloud"
+backendUrl: "http://localhost:8000"
```

<details>
<summary>Frontend running in Docker</summary>

If you run your frontend in a Docker network,
but the `gotestproxy` instance is in the host network,
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
-backendUrl: "https://real-backend-url.in-the.cloud"
+backendUrl: "http://192.168.1.101:8000"
```

and run `./gotestproxy -target=https://real-backend-url.in-the.cloud -port 8000 -color`.
</details>

<details>
<summary>Setup for a backend client</summary>

When setting up `gotestproxy` for the backend (the backend is a client to another service),
look for the URL of the service you want to proxy requests to.
You can set up many `gotestproxy` instances for different services, for example:

```diff
--- a/application.yaml
+++ b/application.yaml
 services:
-  routingService: "https://real-routing-service.in-the.cloud"
-  invoiceService: "https://real-invoice-service.in-the.cloud"
+  routingService: "http://localhost:8000"
+  invoiceService: "http://localhost:8001"
```

and run multiple `./gotestproxy` instances on ports `8000` and `8001`.
</details>

## Mocking backend responses

See [docs/mocking.md](docs/mocking.md).

Then rebuild the project and run it with `-allMocks` arg to enable all mock groups,
or `-mockGroups=group1,group2` arg to enable specified mock groups.

## Using monitors

The monitors let you print information to the console whenever a request or response goes through `gotestproxy`.

### Available monitors

- `request` / `response` monitors: shows basic HTTP info: method, headers, body, status
- `curl` monitor: shows the request as a curl command that you can copy and paste into another program
  like Bruno or Insomnia, or run directly in a shell

### Adjusting which monitors are enabled

Run `gotestproxy` with the `monitors` parameter with the monitors you want to set,
e.g., `-monitors=request,response` if you only want to see request and response, but not curl.

The order in the `monitors` parameter will decide what gets shown first:
`-monitors=request,curl` and `-monitors=curl,request` will produce different output.

### Writing new monitors

See [docs/implementing-monitor.md](docs/implementing-monitor.md).

## Install

When mocking endpoints, installing `gotestproxy` on the system is usually not very handy,
but if you need to have `gotestproxy` available in your `PATH`,
see [docs/installing.md](docs/installing.md)

## Other details

<details>
<summary>Building for other targets</summary>

You can build `gotestproxy` for different targets, like windows/amd64:

```sh
make build-windows-amd64
```

Look at the `MAKEFILE` for examples.
</details>

<details>
<summary>Why Go?</summary>

I used Go because it allowed me to avoid any external dependencies
and rely only on the standard library.
This in turn made this project more plug-and-play,
easier to set up, more portable, and very resistant to supply chain attacks.

The Go compiler is also very fast, so rebuilding the project takes almost no time.
</details>
