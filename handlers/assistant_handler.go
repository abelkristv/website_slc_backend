package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	slcPositionStr := r.URL.Query().Get("slcposition")
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
	} else {
		users, err = h.assistantService.GetAllAssistants()
		if err != nil {
			http.Error(w, "Unable to retrieve assistants", http.StatusInternalServerError)
			return
		}
	}

	if name != "" {
		filteredUsers := []models.Assistant{}
		for _, user := range users {
			userInitialGeneration := user.Initial + user.Generation
			if contains(user.FullName, name) || contains(userInitialGeneration, name) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
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

	if slcPositionStr != "" {
		slcPosition, err := strconv.Atoi(slcPositionStr)
		if err == nil {
			filteredUsers := []models.Assistant{}
			for _, user := range users {
				if user.SLCPositionID == uint(slcPosition) {
					filteredUsers = append(filteredUsers, user)
				}
			}
			users = filteredUsers
		}
	}

	// Calculate total count before pagination
	totalCount := len(users)

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
	if startIndex > totalCount {
		users = []models.Assistant{}
	} else if endIndex > totalCount {
		users = users[startIndex:]
	} else {
		users = users[startIndex:endIndex]
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit // This handles any remainder

	// Create a response structure
	response := struct {
		Users      []models.Assistant `json:"users"`
		TotalCount int                `json:"total_count"`
		TotalPages int                `json:"total_pages"`
	}{
		Users:      users,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		return
	}

	// log.Print(user)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *AssistantHandler) CreateAssistant(w http.ResponseWriter, r *http.Request) {
	var assistant models.Assistant

	if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// log.Print(assistant)

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
	// log.Print(r)
	log.Print(assistant)

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
