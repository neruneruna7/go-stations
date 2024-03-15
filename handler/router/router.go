package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	var healthzHandler = middleware.Middlwares(handler.NewHealthzHandler())
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)

	// Todoについて
	var TODOService = service.NewTODOService(todoDB)
	var todoHandler = middleware.Middlwares(handler.NewTODOHandler(TODOService))
	mux.Handle("/todos", todoHandler)

	// 必ずpanicを起こす

	// var doPanicHandler = middleware.Recovery(middleware.CaptureDeviceOs(handler.NewDoPanicHandler()))
	// var doPanicHandler = middleware.Recovery(handler.NewDoPanicHandler())
	var doPanicHandler = middleware.Middlwares(handler.NewDoPanicHandler())
	mux.Handle("/do-panic", doPanicHandler)
	return mux
}
