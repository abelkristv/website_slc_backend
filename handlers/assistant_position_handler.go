package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AssistantPositionHandler struct {
	service services.AssistantPositionService
}

func NewAssistantPositionHandler(service services.AssistantPositionService) *AssistantPositionHandler {
	return &AssistantPositionHandler{
		service: service,
	}
}
func (h *AssistantPositionHandler) CreatePositionByAssistant(w http.ResponseWriter, r *http.Request) {
	var request struct {
		AssistantId  int       `json:"assistant_id"`
		PositionName string    `json:"position_name"`
		PositionDesc string    `json:"position_desc"`
		StartDate    time.Time `json:"start_date"`
		EndDate      time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.service.CreatePositionByAssistant(
		request.AssistantId,
		request.PositionName,
		request.PositionDesc,
		request.StartDate,
		request.EndDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "Position created successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *AssistantPositionHandler) GetAllAssistantPositions(w http.ResponseWriter, r *http.Request) {
	assistantPositions, err := h.service.GetAllAssistantPositions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assistantPositions)
}

func (h *AssistantPositionHandler) GetAssistantPositionById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	assistantPosition, err := h.service.GetAssistantPositionById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if assistantPosition == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assistantPosition)
}

func (h *AssistantPositionHandler) UpdateAssistantPosition(w http.ResponseWriter, r *http.Request) {
	var request models.AssistantPosition
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAssistantPosition(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AssistantPositionHandler) DeleteAssistantPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	assistantPosition := &models.AssistantPosition{Model: gorm.Model{ID: uint(id)}}
	if err := h.service.DeleteAssistantPosition(assistantPosition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
