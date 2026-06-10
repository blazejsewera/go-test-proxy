package main

import (
	"context"
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
	err := server.ListenAndServe()
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
