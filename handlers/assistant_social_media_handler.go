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

type AssistantSocialMediaHandler struct {
	service services.AssistantSocialMediaService
}

func NewAssistantSocialMediaHandler(service services.AssistantSocialMediaService) *AssistantSocialMediaHandler {
	return &AssistantSocialMediaHandler{service}
}

func (h *AssistantSocialMediaHandler) CreateOrUpdateAssistantSocialMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var socialMedia models.AssistantSocialMedia
	if err := json.NewDecoder(r.Body).Decode(&socialMedia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	socialMedia.AssistantId = int(userID)

	existingSocialMedia, err := h.service.FindByAssistantId(socialMedia.AssistantId)
	if err != nil {
		http.Error(w, "Error checking existing social media", http.StatusInternalServerError)
		return
	}

	if existingSocialMedia != nil {
		existingSocialMedia.GithubLink = socialMedia.GithubLink
		existingSocialMedia.LinkedInLink = socialMedia.LinkedInLink
		existingSocialMedia.InstagramLink = socialMedia.InstagramLink
		existingSocialMedia.WhatsappLink = socialMedia.WhatsappLink
		existingSocialMedia.PersonalWebsiteLink = socialMedia.PersonalWebsiteLink

		if err := h.service.UpdateAssistantSocialMedia(existingSocialMedia); err != nil {
			http.Error(w, "Error updating social media", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingSocialMedia)
	} else {
		if err := h.service.CreateAssistantSocialMedia(&socialMedia); err != nil {
			http.Error(w, "Error creating social media", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(socialMedia)
	}
}

func (h *AssistantSocialMediaHandler) GetAssistantSocialMediaByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	socialMedia, err := h.service.GetAssistantSocialMediaByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(socialMedia)
}

func (h *AssistantSocialMediaHandler) UpdateAssistantSocialMedia(w http.ResponseWriter, r *http.Request) {
	var socialMedia models.AssistantSocialMedia
	if err := json.NewDecoder(r.Body).Decode(&socialMedia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAssistantSocialMedia(&socialMedia); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(socialMedia)
}

func (h *AssistantSocialMediaHandler) DeleteAssistantSocialMedia(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAssistantSocialMedia(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
