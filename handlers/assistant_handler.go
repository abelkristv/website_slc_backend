package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode"

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

// GetAllAssistants retrieves assistants based on generation and name filters.
func (h *AssistantHandler) GetAllAssistants(w http.ResponseWriter, r *http.Request) {
	generation := r.URL.Query().Get("generation")
	name := r.URL.Query().Get("name") // Get the name from query parameters

	var users []models.Assistant
	var err error

	if generation != "" {
		// If a generation is provided, fetch assistants by generation
		users, err = h.assistantService.GetAssistantsByGeneration(generation)
		if err != nil {
			http.Error(w, "Unable to retrieve assistants for the specified generation", http.StatusInternalServerError)
			return
		}

		// Check the first two letters of the name
		if name != "" {
			if isInitialsSearch(name) { // Check if the first two letters are uppercase
				users = filterByInitials(users, name) // Filter by initials
			} else {
				users = filterByName(users, name) // Filter by name
			}
		}
	} else if name != "" {
		// If no generation is provided but a name is provided, fetch all assistants
		users, err = h.assistantService.GetAllAssistants()
		if err != nil {
			http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
			return
		}

		// Check the first two letters of the name
		if isInitialsSearch(name) {
			users = filterByInitials(users, name) // Filter by initials
		} else {
			users = filterByName(users, name) // Filter by name
		}
	} else {
		// If neither generation nor name is provided, fetch all assistants
		users, err = h.assistantService.GetAllAssistants()
		if err != nil {
			http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// isInitialsSearch checks if the first two characters of the name are uppercase.
func isInitialsSearch(name string) bool {
	if len(name) < 2 {
		return false // Not enough characters to check
	}
	return unicode.IsUpper(rune(name[0])) && unicode.IsUpper(rune(name[1]))
}

// filterByName filters users based on the full name containing the name parameter.
func filterByName(users []models.Assistant, name string) []models.Assistant {
	filteredUsers := []models.Assistant{}
	for _, user := range users {
		if contains(user.FullName, name) { // Check if the user's full name contains the name parameter
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

// filterByInitials filters users based on initials containing the name parameter.
func filterByInitials(users []models.Assistant, name string) []models.Assistant {
	filteredUsers := []models.Assistant{}
	for _, user := range users {
		if strings.Contains(user.Initial, name) { // Check if the user's initials contain the name parameter
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

// Helper function to check if a string contains another string (case insensitive)
func contains(fullName, name string) bool {
	return strings.Contains(strings.ToLower(fullName), strings.ToLower(name))
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

func (h *AssistantHandler) GetAllGenerations(w http.ResponseWriter, r *http.Request) {
	generations, err := h.assistantService.GetAllGenerations()
	if err != nil {
		http.Error(w, "Unable to retrieve generations", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(generations)
}

func (h *AssistantHandler) GetAssistantsByGeneration(w http.ResponseWriter, r *http.Request) {
	generation := r.URL.Query().Get("generation") // Get the generation from query parameters

	if generation == "" {
		http.Error(w, "Generation parameter is required", http.StatusBadRequest)
		return
	}

	assistants, err := h.assistantService.GetAssistantsByGeneration(generation)
	if err != nil {
		http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assistants)
}
