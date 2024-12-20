package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abelkristv/slc_website/services"
)

type EmailHandler struct {
	EmailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{EmailService: emailService}
}

func (h *EmailHandler) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.EmailService.SendEmail(req.To, req.Subject, req.Body); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
