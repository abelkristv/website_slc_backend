// handlers/teaching_history_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
)

type TeachingHistoryHandler struct {
	service services.TeachingHistoryService
}

func NewTeachingHistoryHandler(service services.TeachingHistoryService) *TeachingHistoryHandler {
	return &TeachingHistoryHandler{service: service}
}

func (h *TeachingHistoryHandler) GetTeachingHistoryByAssistantAndPeriod(w http.ResponseWriter, r *http.Request) {
	assistantUsername := strings.ToLower(r.URL.Query().Get("assistant_username"))
	if assistantUsername == "" {
		http.Error(w, "Missing or invalid assistant_username", http.StatusBadRequest)
		return
	}

	periodName := strings.ToLower(r.URL.Query().Get("period_name"))
	if periodName == "" {
		http.Error(w, "Missing or invalid period_name", http.StatusBadRequest)
		return
	}

	histories, err := h.service.GetTeachingHistoryByAssistantAndPeriod(assistantUsername, periodName)
	if err != nil {
		http.Error(w, "Error retrieving teaching history", http.StatusInternalServerError)
		return
	}

	type TeachingHistoryResponse struct {
		ID     uint          `json:"id"`
		Course models.Course `json:"course"`
	}

	var response []TeachingHistoryResponse
	for _, history := range histories {
		response = append(response, TeachingHistoryResponse{
			ID:     history.ID,
			Course: history.Course,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TeachingHistoryHandler) GetTeachingHistoryGroupedByPeriod(w http.ResponseWriter, r *http.Request) {
	assistantUsername := strings.ToLower(r.URL.Query().Get("assistant_username"))
	if assistantUsername == "" {
		http.Error(w, "Missing or invalid assistant_username", http.StatusBadRequest)
		return
	}

	groupedHistories, err := h.service.GetTeachingHistoryGroupedByPeriod(assistantUsername)
	if err != nil {
		http.Error(w, "Error retrieving teaching history", http.StatusInternalServerError)
		return
	}
	type PeriodCoursesResponse struct {
		PeriodName string          `json:"periodName"`
		Courses    []models.Course `json:"courses"`
	}

	// Create the response structure
	var response []PeriodCoursesResponse

	for periodName, histories := range groupedHistories {
		var courses []models.Course

		for _, history := range histories {
			courses = append(courses, history.Course)
		}

		response = append(response, PeriodCoursesResponse{
			PeriodName: periodName,
			Courses:    courses,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
