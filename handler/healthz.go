package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res = &model.HealthzResponse{
		Message: "OK",
	}
	var encoder = json.NewEncoder(w)
	// これResultみたいなものかと思ったら，エラーが返ってくるのか
	var err = encoder.Encode(res)
	if err != nil {
		log.Println(err)
	}
}
