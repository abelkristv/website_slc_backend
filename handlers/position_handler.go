package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type PositionHandler struct {
	positionService services.PositionService
}

func NewPositionHandler(positionService services.PositionService) *PositionHandler {
	return &PositionHandler{
		positionService: positionService,
	}
}

func (h *PositionHandler) GetAllPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := h.positionService.GetAllPositions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(positions)
}

func (h *PositionHandler) GetPositionById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	position, err := h.positionService.GetPositionById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if position == nil {
		http.Error(w, "Position not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(position)
}

func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var position models.Position
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.positionService.CreatePosition(&position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(position)
}

func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	var position models.Position
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	position.ID = uint(id)

	err = h.positionService.UpdatePosition(&position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(position)
}

func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	err = h.positionService.DeletePosition(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
