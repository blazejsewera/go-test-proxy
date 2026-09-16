package main

import (
	"context"
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

	consoleMonitor := monitor.NewConsoleMonitor(cfg.Target, cfg.Cfmt)
	curlRequestMonitor := monitor.NewCurlRequestMonitor(cfg.Target, cfg.Cfmt)
	stderrMonitor := monitor.NewStdErrMonitor(cfg.Cfmt)

	builder := proxy.NewBuilder().
		WithProxyTarget(cfg.Target).
		WithPort(cfg.Port).
		WithMockGroup(mock.ConfigMockGroup, mock.ConfigInfo(cfg)).
		WithMonitor(monitor.Combine(consoleMonitor, curlRequestMonitor, stderrMonitor))

	if cfg.EnableAllMocks {
		builder.WithEnabledAllMocks()
	} else {
		builder.WithEnabledMockGroups(cfg.MockGroups)
	}

	server := builder.Build()

	listenAndServe(server, cfg)
	defer shutdownServer(server)
}

func listenAndServe(server *http.Server, cfg config.Configuration) {
	log.Printf("starting proxy server for target: '%s', on port: %d", cfg.Target, cfg.Port)
	network := "tcp4"
	if cfg.IPv6 {
		log.Printf(", on IPv6")
		network = "tcp6"
	}
	l, err := net.Listen(network, fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	err = server.Serve(l)
	if err != nil {
		log.Fatalln(err)
	}
}

func shutdownServer(server *http.Server) {
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
