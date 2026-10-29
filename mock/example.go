package mock

import (
	"encoding/json"
	"net/http"

	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
)

const (
	ExampleRoute = "/_example"
)

type response struct {
	ResponseValue string `json:"responseValue"`
}

func ExampleHandler200(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParam := r.URL.Query().Get("queryParam")
	res := response{queryParam}

	responseJSON, _ := json.Marshal(res)
	_, err := w.Write(responseJSON)
	if err != nil {
		log.Fatalf("error in ExampleHandler: %v\n", err)
	}
}

func ExampleHandler404(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
