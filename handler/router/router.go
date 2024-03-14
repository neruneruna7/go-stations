package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	var healthzHandler = handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)

	var TODOService = service.NewTODOService(todoDB)
	var todoHandler = handler.NewTODOHandler(TODOService)
	mux.Handle("/todos", todoHandler)
	return mux
}
