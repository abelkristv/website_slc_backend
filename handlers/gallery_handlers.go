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

func (h *GalleryHandler) GetPendingGalleries(w http.ResponseWriter, r *http.Request) {
	galleries, err := h.service.GetGalleriesByStatus("pending")
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

    if gallery.AssistantId != user.AssistantId {
        http.Error(w, "Unauthorized to update this gallery", http.StatusForbidden)
        return
    }

    if err := h.service.UpdateGallery(&gallery); err != nil {
        http.Error(w, "Failed to update gallery", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(gallery)
}

func (h *GalleryHandler) DeleteGallery(w http.ResponseWriter, r *http.Request) {
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

    galleryID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    gallery, err := h.service.GetGalleryByID(uint(galleryID))
    if err != nil {
        http.Error(w, "Gallery not found", http.StatusNotFound)
        return
    }

    if gallery.AssistantId != user.AssistantId {
        http.Error(w, "Unauthorized to delete this gallery", http.StatusForbidden)
        return
    }

    if err := h.service.DeleteGallery(uint(galleryID)); err != nil {
        http.Error(w, "Failed to delete gallery", http.StatusInternalServerError)
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