package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB, basic_auth_config *model.BasicAuthConfig) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	var healthzHandler = middleware.CommonMiddlwares(handler.NewHealthzHandler())
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)

	// Todoについて
	var TODOService = service.NewTODOService(todoDB)
	var todoHandler = middleware.CommonMiddlwares(handler.NewTODOHandler(TODOService))
	mux.Handle("/todos", todoHandler)

	// 必ずpanicを起こす

	// var doPanicHandler = middleware.Recovery(middleware.CaptureDeviceOs(handler.NewDoPanicHandler()))
	// var doPanicHandler = middleware.Recovery(handler.NewDoPanicHandler())
	var doPanicHandler = middleware.CommonMiddlwares(middleware.BasicAuth(basic_auth_config, handler.NewDoPanicHandler()))
	mux.Handle("/do-panic", doPanicHandler)
	return mux
}
