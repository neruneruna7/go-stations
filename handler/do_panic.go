package handler

import (
	"net/http"
)

type DoPanicHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewDoPanicHandler() *DoPanicHandler {
	return &DoPanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *DoPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Oops! I destroyed the railroad tracks!")
}
