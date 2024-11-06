package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abelkristv/slc_website/middleware"
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Unable to retrieve users", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newUser, err := h.userService.CreateUser(user.Username, user.Password, user.Role, user.AssistantId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.ID = uint(id)
	if err := h.userService.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := strings.ToUpper(credentials.Username)
	token, userToReturn, err := h.userService.Login(username, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"data": map[string]interface{}{
			"token": token,
			"user":  userToReturn,
		},
	})

}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Print("User ID:", userID)

	// Fetch the user information from the service layer
	user, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return the user data as JSON
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Set the cookie's expiration date to the past to invalidate it
	http.SetCookie(w, &http.Cookie{
		Name:     "token",         // Name of the cookie
		Value:    "",              // Clear the token
		Path:     "/",             // Cookie path
		Expires:  time.Unix(0, 0), // Set expiration date to the past
		HttpOnly: true,            // Make it HTTP-only
		Secure:   true,            // Use secure flag if you're using HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Respond with a message indicating successful logout
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out"))
}
