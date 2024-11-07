package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type AssistantAwardHandler struct {
	service services.AssistantAwardService
}

func NewAssistantAwardHandler(service services.AssistantAwardService) *AssistantAwardHandler {
	return &AssistantAwardHandler{service}
}

func (h *AssistantAwardHandler) CreateAssistantAward(w http.ResponseWriter, r *http.Request) {
	var assistantAward models.AssistantAward
	if err := json.NewDecoder(r.Body).Decode(&assistantAward); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateAssistantAward(&assistantAward); err != nil {
		http.Error(w, "Could not create assistant award", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assistantAward)
}

func (h *AssistantAwardHandler) GetAssistantAwardByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid award ID", http.StatusBadRequest)
		return
	}

	assistantAward, err := h.service.GetAssistantAwardByID(uint(id))
	if err != nil {
		http.Error(w, "Assistant award not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(assistantAward)
}

func (h *AssistantAwardHandler) GetAssistantAwardsByAssistantID(w http.ResponseWriter, r *http.Request) {
	assistantIDParam := mux.Vars(r)["assistantId"]
	assistantID, err := strconv.Atoi(assistantIDParam)
	if err != nil {
		http.Error(w, "Invalid assistant ID", http.StatusBadRequest)
		return
	}

	assistantAwards, err := h.service.GetAssistantAwardsByAssistantID(assistantID)
	if err != nil {
		http.Error(w, "Could not retrieve assistant awards", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assistantAwards)
}

func (h *AssistantAwardHandler) UpdateAssistantAward(w http.ResponseWriter, r *http.Request) {
	var assistantAward models.AssistantAward
	if err := json.NewDecoder(r.Body).Decode(&assistantAward); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAssistantAward(&assistantAward); err != nil {
		http.Error(w, "Could not update assistant award", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assistantAward)
}

func (h *AssistantAwardHandler) DeleteAssistantAward(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid award ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAssistantAward(uint(id)); err != nil {
		http.Error(w, "Could not delete assistant award", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
