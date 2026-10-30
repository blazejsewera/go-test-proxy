package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/mock"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy"
)

func main() {
	cfg := config.ParseConfig()

	requestMonitor := monitor.NewRequestMonitor(cfg.Target, cfg.Cfmt)
	responseMonitor := monitor.NewResponseMonitor(cfg.Target, cfg.Cfmt)
	curlRequestMonitor := monitor.NewCurlRequestMonitor(cfg.Target, cfg.Cfmt)

	builder := proxy.NewBuilder().
		WithProxyTarget(cfg.Target).
		WithPort(cfg.Port).
		WithMockGroup("example200", proxy.Mock{RoutePattern: mock.ExampleRoute, HandlerFunc: mock.ExampleHandler200}).
		WithMockGroup("example404", proxy.Mock{RoutePattern: mock.ExampleRoute, HandlerFunc: mock.ExampleHandler404}).
		WithMonitor(monitor.Combine(cfg, requestMonitor, responseMonitor, curlRequestMonitor))

	if cfg.EnableAllMocks {
		builder.WithEnabledAllMocks()
	} else {
		builder.WithEnabledMockGroups(cfg.MockGroups)
	}

	router := builder.Build()
	listenAndServe(router, cfg)
}

func listenAndServe(router *http.ServeMux, cfg config.Configuration) {
	log.Printf("starting proxy server for target: '%s', on port: %d", cfg.Target, cfg.Port)
	network := "tcp4"
	if cfg.IPv6 {
		log.Printf(", on IPv6")
		network = "tcp6"
	}
	log.Printf("\n")

	l, err := net.Listen(network, fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer shutdownListener(l)
	err = http.Serve(l, router)
	if err != nil {
		log.Fatalln(err)
	}
}

func shutdownListener(l net.Listener) {
	err := l.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
