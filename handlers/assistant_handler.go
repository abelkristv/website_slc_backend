package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type AssistantHandler struct {
	assistantService *services.AssistantService
}

func NewAssistantHandler(assistantService *services.AssistantService) *AssistantHandler {
	return &AssistantHandler{
		assistantService: assistantService,
	}
}

func (h *AssistantHandler) GetAllAssistants(w http.ResponseWriter, r *http.Request) {
	users, err := h.assistantService.GetAllAssistants()
	if err != nil {
		http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(users)
}

func (h *AssistantHandler) GetAssistantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid assistant ID", http.StatusBadRequest)
		return
	}
	user, err := h.assistantService.GetAssistantById(uint(id))
	if err != nil {
		http.Error(w, "Assistant not found", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(user)
}

func (h *AssistantHandler) CreateAssistant(w http.ResponseWriter, r *http.Request) {
	var assistant models.Assistant

	if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newAssistant, err := h.assistantService.CreateAssistant(assistant.Email, assistant.Bio, assistant.ProfilePicture, assistant.Initial, assistant.Generation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAssistant)
}

func (h *AssistantHandler) UpdateAssistant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid assistant ID", http.StatusBadRequest)
		return
	}

	var assistant models.Assistant
	if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	assistant.ID = uint(id)
	if err := h.assistantService.UpdateAssistant(&assistant); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AssistantHandler) DeleteAssistant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.assistantService.DeleteAssistant(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
