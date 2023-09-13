package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/blazejsewera/go-test-proxy/monitor"
	"github.com/blazejsewera/go-test-proxy/proxy"
	"log"
	"net/http"
)

func main() {
	info := ConfigInfo{Application: "https://github.com/blazejsewera/go-test-proxy"}
	builder := proxy.NewBuilder()

	target := flag.String("target", "", "the target host address, for example, http://example.com")
	port := flag.Int("port", 8000, "the port on which the proxy server will be running")
	flag.Parse()

	if target != nil && *target != "" {
		info.Target = *target
		builder.WithProxyTarget(*target)
	}

	if port != nil && *port != 0 {
		portUint16 := uint16(*port)
		info.Port = portUint16
		builder.WithPort(portUint16)
	}

	builder.WithHandlerFunc("/_info", configInfoHandler(info))

	consoleMonitor := monitor.NewConsoleMonitor(info.Target)
	curlRequestMonitor := monitor.NewCurlRequestMonitor(info.Target)
	stderrMonitor := monitor.NewStdErrMonitor()
	builder.WithMonitor(monitor.Combine(consoleMonitor, curlRequestMonitor, stderrMonitor))

	server := builder.Build()
	log.Printf("starting proxy server for target: '%s', go to 'http://localhost:%d/_info' to get config", info.Target, info.Port)
	listenAndServe(server)
	defer shutdownServer(server)
}

func listenAndServe(server *http.Server) {
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

func configInfoHandler(info ConfigInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := json.Marshal(info)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(bodyBytes)
		if err != nil {
			return
		}
	}
}

type ConfigInfo struct {
	Application string `json:"application"`
	Port        uint16 `json:"port"`
	Target      string `json:"target"`
}
