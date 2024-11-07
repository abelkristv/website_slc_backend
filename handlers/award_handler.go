// handlers/award_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type AwardHandler struct {
	service services.AwardService
}

func NewAwardHandler(service services.AwardService) *AwardHandler {
	return &AwardHandler{service}
}

func (h *AwardHandler) CreateAward(w http.ResponseWriter, r *http.Request) {
	var award models.Award
	if err := json.NewDecoder(r.Body).Decode(&award); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateAward(&award); err != nil {
		http.Error(w, "Could not create award", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(award)
}

func (h *AwardHandler) GetAwardByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid award ID", http.StatusBadRequest)
		return
	}

	award, err := h.service.GetAwardByID(uint(id))
	if err != nil {
		http.Error(w, "Award not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(award)
}

func (h *AwardHandler) UpdateAward(w http.ResponseWriter, r *http.Request) {
	var award models.Award
	if err := json.NewDecoder(r.Body).Decode(&award); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAward(&award); err != nil {
		http.Error(w, "Could not update award", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(award)
}

func (h *AwardHandler) DeleteAward(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid award ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAward(uint(id)); err != nil {
		http.Error(w, "Could not delete award", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AwardHandler) GetAllAwards(w http.ResponseWriter, r *http.Request) {
	awards, err := h.service.GetAllAwards()
	if err != nil {
		http.Error(w, "Could not retrieve awards", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(awards)
}
