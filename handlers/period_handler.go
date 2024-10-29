package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type PeriodHandler struct {
	periodService *services.PeriodService
}

func NewPeriodHandler(periodService *services.PeriodService) *PeriodHandler {
	return &PeriodHandler{
		periodService: periodService,
	}
}

func (h *PeriodHandler) GetAllPeriods(w http.ResponseWriter, r *http.Request) {
	periods, err := h.periodService.GetAllPeriods()
	if err != nil {
		http.Error(w, "Unable to retrieve events", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(periods)
}

func (h *PeriodHandler) GetPeriodById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid period ID", http.StatusBadRequest)
		return
	}
	event, err := h.periodService.GetPeriodById(uint(id))
	if err != nil {
		http.Error(w, "Period not found", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(event)
}

func (h *PeriodHandler) CreatePeriod(w http.ResponseWriter, r *http.Request) {
	var period models.Period

	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newEvent, err := h.periodService.CreatePeriod(period.PeriodTitle, period.StartDate, period.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (h *PeriodHandler) UpdatePeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid period ID", http.StatusBadRequest)
		return
	}

	var period models.Period
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	period.ID = uint(id)
	if err := h.periodService.UpdatePeriod(&period); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PeriodHandler) DeletePeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid period ID", http.StatusBadRequest)
		return
	}

	if err := h.periodService.DeletePeriod(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
