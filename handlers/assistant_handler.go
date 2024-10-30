package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
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

func (h *AssistantHandler) GetAllAssistants(w http.ResponseWriter, r *http.Request) {
	generation := strings.ToLower(r.URL.Query().Get("generation"))
	name := strings.ToLower(r.URL.Query().Get("name"))
	orderby := strings.ToLower(r.URL.Query().Get("orderby"))
	order := strings.ToLower(r.URL.Query().Get("order"))
	status := strings.ToLower(r.URL.Query().Get("status"))
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var users []models.Assistant
	var err error

	page := 1
	limit := 24
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if generation != "" {
		users, err = h.assistantService.GetAssistantsByGeneration(generation)
		if err != nil {
			http.Error(w, "Unable to retrieve assistants for the specified generation", http.StatusInternalServerError)
			return
		}

		if name != "" {
			if isInitialsSearch(name) {
				users = filterByInitials(users, name)
			} else {
				users = filterByName(users, name)
			}
		}
	} else if name != "" {
		users, err = h.assistantService.GetAllAssistants()
		if err != nil {
			http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
			return
		}

		if isInitialsSearch(name) {
			users = filterByInitials(users, name)
		} else {
			users = filterByName(users, name)
		}
	} else {
		users, err = h.assistantService.GetAllAssistants()
		if err != nil {
			http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
			return
		}
	}

	if status != "" {
		filteredUsers := []models.Assistant{}
		for _, user := range users {
			if (status == "active" && user.Status == "active") || (status == "inactive" && user.Status == "inactive") {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	if orderby != "" {
		sort.Slice(users, func(i, j int) bool {
			var less bool
			switch orderby {
			case "generation":
				less = users[i].Generation < users[j].Generation
			case "name":
				less = users[i].FullName < users[j].FullName
			case "initial":
				less = users[i].Initial < users[j].Initial
			default:
				less = true
			}

			if order == "descending" {
				return !less
			}
			return less
		})
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if startIndex > len(users) {
		users = []models.Assistant{}
	} else if endIndex > len(users) {
		users = users[startIndex:]
	} else {
		users = users[startIndex:endIndex]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func isInitialsSearch(name string) bool {
	if len(name) < 2 {
		return false
	}
	return unicode.IsUpper(rune(name[0])) && unicode.IsUpper(rune(name[1]))
}

func filterByName(users []models.Assistant, name string) []models.Assistant {
	filteredUsers := []models.Assistant{}
	for _, user := range users {
		if contains(user.FullName, name) { // Check if the user's full name contains the name parameter
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func filterByInitials(users []models.Assistant, name string) []models.Assistant {
	filteredUsers := []models.Assistant{}
	for _, user := range users {
		if strings.Contains(user.Initial, name) { // Check if the user's initials contain the name parameter
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

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
