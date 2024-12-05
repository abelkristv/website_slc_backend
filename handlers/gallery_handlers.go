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

type GalleryHandler struct {
	service     services.GalleryService
	userService services.UserService
}

func NewGalleryHandler(service services.GalleryService, userService services.UserService) *GalleryHandler {
	return &GalleryHandler{
		service:     service,
		userService: userService,
	}
}

func (h *GalleryHandler) CreateGallery(w http.ResponseWriter, r *http.Request) {
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

	var gallery models.Gallery
	if err := json.NewDecoder(r.Body).Decode(&gallery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(gallery.GalleryImages) == 0 {
		http.Error(w, "No images provided", http.StatusBadRequest)
		return
	}

	gallery.AssistantId = user.AssistantId

	if user.Assistant.SLCPosition.PositionName == "Operations Management Officer" {
		gallery.GalleryStatus = "accepted"
	} else {
		gallery.GalleryStatus = "pending"
	}

	if err := h.service.CreateGallery(&gallery); err != nil {
		http.Error(w, "Failed to create gallery", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gallery)
}

func (h *GalleryHandler) GetAcceptedGalleries(w http.ResponseWriter, r *http.Request) {
	galleries, err := h.service.GetGalleriesByStatus("accepted")
	if err != nil {
		http.Error(w, "Failed to retrieve galleries", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(galleries)
}

func (h *GalleryHandler) GetAllGalleries(w http.ResponseWriter, r *http.Request) {
	galleries, err := h.service.GetAllGalleries()
	if err != nil {
		http.Error(w, "Failed to retrieve galleries", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(galleries)
}

func (h *GalleryHandler) GetMyGalleries(w http.ResponseWriter, r *http.Request) {
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

	galleries, err := h.service.GetGalleriesByAssistantID(uint(user.AssistantId))
	if err != nil {
		http.Error(w, "Failed to retrieve galleries", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(galleries)
}

func (h *GalleryHandler) UpdateGallery(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updatedGallery models.Gallery
	if err := json.NewDecoder(r.Body).Decode(&updatedGallery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	galleryID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return
	}
	updatedGallery.ID = uint(galleryID)

	err = h.service.UpdateGallery(userID, &updatedGallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Gallery updated successfully"})
}

func (h *GalleryHandler) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	galleryID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteGallery(userID, uint(galleryID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Gallery deleted successfully"})
}

func (h *GalleryHandler) AcceptGallery(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return
	}

	gallery, err := h.service.GetGalleryByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve gallery", http.StatusInternalServerError)
		return
	}

	err = h.service.AcceptGallery(gallery)
	if err != nil {
		http.Error(w, "Failed to accept gallery", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gallery accepted successfully"))
}

func (h *GalleryHandler) RejectGallery(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return
	}

	gallery, err := h.service.GetGalleryByID(uint(id))
	if err != nil {
		http.Error(w, "Failed to retrieve gallery", http.StatusInternalServerError)
		return
	}

	err = h.service.RejectGallery(gallery)
	if err != nil {
		http.Error(w, "Failed to reject gallery", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gallery rejected successfully"))
}
