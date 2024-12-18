package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type ContactUsHandler struct {
	service services.ContactUsService
	emailService   services.EmailService
}

func NewContactUsHandler(contactService services.ContactUsService, emailService services.EmailService) *ContactUsHandler {
	return &ContactUsHandler{
		service: contactService,
		emailService:   emailService,
	}
}

func (h *ContactUsHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.service.GetAllContacts()
	if err != nil {
		http.Error(w, "Error fetching contacts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactUsHandler) GetContactById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := h.service.GetContactById(uint(id))
	if err != nil {
		http.Error(w, "Error fetching contact", http.StatusInternalServerError)
		return
	}
	if contact == nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactUsHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.ContactUs

	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateContact(&contact); err != nil {
		http.Error(w, "Error creating contact", http.StatusInternalServerError)
		return
	}

	subject := "New Contact Form Submission"
	body := "You have received a new message from the contact form:\n\n" +
		"Name: " + contact.Name + "\n" +
		"Email: " + contact.Email + "\n" +
		"Phone: " + contact.Phone + "\n" +
		"Message:\n" + contact.Message + "\n"

	err := h.emailService.SendEmail("jason.tanuarta@binus.ac.id", subject, body)
	if err != nil {
		http.Error(w, "Contact created but failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactUsHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	var contact models.ContactUs
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	contact.ID = uint(id)

	if err := h.service.UpdateContact(&contact); err != nil {
		http.Error(w, "Error updating contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}

func (h *ContactUsHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteContact(uint(id)); err != nil {
		http.Error(w, "Error deleting contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ContactUsHandler) UpdateIsRead(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkContactAsRead(uint(id), true); err != nil {
		http.Error(w, "Error updating is_read status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "updating isRead attribute successfull",
	})
}
