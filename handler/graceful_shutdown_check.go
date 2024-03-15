package handler

import (
	"net/http"
	"time"
)

type GracefulShutdownHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewGracefulShutdownHandler() *GracefulShutdownHandler {
	return &GracefulShutdownHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *GracefulShutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 5秒待ってからokを返す
	// これはGraceful Shutdownのテスト用のハンドラ
	time.Sleep(time.Second * 7)
	http.Error(w, "OK", http.StatusOK)
}
