package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type SLCPositionHandler struct {
	service services.SLCPositionService
}

func NewSLCPositionHandler(service services.SLCPositionService) *SLCPositionHandler {
	return &SLCPositionHandler{service: service}
}

func (h *SLCPositionHandler) GetAllPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := h.service.GetAllPositions()
	if err != nil {
		http.Error(w, "Failed to fetch positions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(positions)
}

func (h *SLCPositionHandler) GetPositionByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	position, err := h.service.GetPositionByID(uint(id))
	if err != nil {
		http.Error(w, "Position not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(position)
}

func (h *SLCPositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var position models.SLCPosition
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePosition(&position); err != nil {
		http.Error(w, "Failed to create position", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(position)
}

func (h *SLCPositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var position models.SLCPosition
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePosition(uint(id), &position); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(position)
}

func (h *SLCPositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePosition(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Position deleted successfully"})
}
