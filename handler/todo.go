package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	var todo, e = h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if e != nil {
		return nil, e
	}

	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.TodoPostHandler(w, r)
	case "PUT":
		h.TodoPutHandler(w, r)
	default:
	}
}

func CreateTodoRequestDecode(req *model.CreateTODORequest, r *http.Request) error {
	var decoder = json.NewDecoder(r.Body)
	return decoder.Decode(&req)
}

func CreateTodoResponseEncode(res *model.CreateTODOResponse, w http.ResponseWriter) error {
	var encoder = json.NewEncoder(w)
	return encoder.Encode(res)
}

func UpdateTodoRequestDecode(req *model.UpdateTODORequest, r *http.Request) error {
	var decoder = json.NewDecoder(r.Body)
	return decoder.Decode(&req)
}

func UpdateTodoResponseEncode(res *model.UpdateTODOResponse, w http.ResponseWriter) error {
	var encoder = json.NewEncoder(w)
	return encoder.Encode(res)
}

func (h *TODOHandler) TodoPostHandler(w http.ResponseWriter, r *http.Request) {
	// エラーを処理する責務を持つ
	var req model.CreateTODORequest

	err := CreateTodoRequestDecode(&req, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create TODO", http.StatusBadRequest)
		return
	}

	var res = model.CreateTODOResponse{
		TODO: *todo,
	}
	var err2 = CreateTodoResponseEncode(&res, w)
	if err2 != nil {
		log.Println(err2)
		http.Error(w, "failed to encode response", http.StatusBadRequest)
		return
	}
}

func (h *TODOHandler) TodoPutHandler(w http.ResponseWriter, r *http.Request) {
	// エラーを処理する責務を持つ
	var req model.UpdateTODORequest

	err := UpdateTodoRequestDecode(&req, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to json decode", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to update TODO", http.StatusBadRequest)
		return
	}

	var res = model.UpdateTODOResponse{
		TODO: *todo,
	}
	var err2 = UpdateTodoResponseEncode(&res, w)
	if err2 != nil {
		log.Println(err2)
		http.Error(w, "failed to json encode", http.StatusBadRequest)
		return
	}
}
