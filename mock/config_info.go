package mock

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
	"github.com/blazejsewera/go-test-proxy/config"
	"github.com/blazejsewera/go-test-proxy/proxy"
)

const (
	ConfigMockGroup = "config"
	configRoute     = "/_config"
)

func configInfoHandler(cfg config.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := json.Marshal(cfg)
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

func ConfigInfo(cfg config.Configuration) proxy.Mock {
	if cfg.EnableAllMocks || slices.Contains(cfg.MockGroups, ConfigMockGroup) {
		log.Printf("go to 'http://localhost:%d/_config' to get config\n", cfg.Port)
	}
	return proxy.Mock{RoutePattern: configRoute, HandlerFunc: configInfoHandler(cfg)}
}
