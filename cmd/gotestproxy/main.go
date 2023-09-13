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
	config := parseConfig()

	consoleMonitor := monitor.NewConsoleMonitor(config.Target)
	curlRequestMonitor := monitor.NewCurlRequestMonitor(config.Target)
	stderrMonitor := monitor.NewStdErrMonitor()

	server := proxy.NewBuilder().
		WithProxyTarget(config.Target).
		WithPort(config.Port).
		WithHandlerFunc("/_config", configInfoHandler(config)).
		WithMonitor(monitor.Combine(consoleMonitor, curlRequestMonitor, stderrMonitor)).
		Build()

	listenAndServe(server, config)
	defer shutdownServer(server)
}

func parseConfig() Configuration {
	config := Configuration{Application: "https://github.com/blazejsewera/go-test-proxy"}

	target := flag.String("target", "", "the target host address, for example, https://example.com")
	port := flag.Int("port", 8000, "the port on which the proxy server will be running")
	flag.Parse()

	if *target == "" {
		log.Fatalln("[FATAL] The target cannot be empty. " +
			"Specify the target server, e.g., https://example.com\n" +
			"Try: ./gotestproxy --target=https://example.com\n" +
			"Run with -h for help.")
	}

	config.Target = *target

	portUint16 := uint16(*port)
	config.Port = portUint16

	return config
}

func listenAndServe(server *http.Server, config Configuration) {
	log.Printf("starting proxy server for target: '%s', go to 'http://localhost:%d/_config' to get config", config.Target, config.Port)
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

func configInfoHandler(config Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := json.Marshal(config)
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

type Configuration struct {
	Application string `json:"application"`
	Port        uint16 `json:"port"`
	Target      string `json:"target"`
}
