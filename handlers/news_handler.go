package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/middleware"
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type NewsHandler struct {
	service services.NewsService
	userService services.UserService
}

func NewNewsHandler(service services.NewsService, userService services.UserService) *NewsHandler {
	return &NewsHandler{
		service: service,
		userService: userService,
	}
}

func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := h.userService.GetUserByID(userID)
    if err != nil {
        http.Error(w, "Error retrieving user", http.StatusInternalServerError)
        return
    }

    if user.AssistantId == 0 {
        http.Error(w, "No associated Assistant ID for user", http.StatusNotFound)
        return
    }

    var news models.News
    if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(news.NewsImages) == 0 {
        http.Error(w, "No images provided", http.StatusBadRequest)
        return
    }

    news.AssistantId = user.AssistantId

    if err := h.service.CreateNews(&news); err != nil {
        http.Error(w, "Failed to create news", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(news)
}

func (h *NewsHandler) GetNewsByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	news, err := h.service.GetNewsByID(uint(id))
	if err != nil {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func (h *NewsHandler) GetAllNews(w http.ResponseWriter, r *http.Request) {
	news, err := h.service.GetAllNews()
	if err != nil {
		http.Error(w, "Failed to retrieve news", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func (h *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := h.userService.GetUserByID(userID)
    if err != nil {
        http.Error(w, "Error retrieving user", http.StatusInternalServerError)
        return
    }

    if user.AssistantId == 0 {
        http.Error(w, "No associated Assistant ID for user", http.StatusNotFound)
        return
    }

    var news models.News
    if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(news.NewsImages) == 0 {
        http.Error(w, "No images provided", http.StatusBadRequest)
        return
    }

    news.AssistantId = user.AssistantId

    if err := h.service.UpdateNews(&news); err != nil {
        http.Error(w, "Failed to update news", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(news)
}


func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNews(uint(id)); err != nil {
		http.Error(w, "Failed to delete news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "News deleted successfully"})
}
