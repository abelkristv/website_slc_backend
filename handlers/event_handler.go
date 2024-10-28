package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetAllEvents()
	if err != nil {
		http.Error(w, "Unable to retrieve events", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetEventById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	event, err := h.eventService.GetEventById(uint(id))
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newEvent, err := h.eventService.CreateEvent(event.EventTitle, event.EventDescription, event.WriterId, event.Type, event.PeriodId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	event.ID = uint(id)
	if err := h.eventService.UpdateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if err := h.eventService.DeleteEvent(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
