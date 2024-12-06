package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type WiredSyncHandler struct{}

func NewWiredSyncHandler() *WiredSyncHandler {
	return &WiredSyncHandler{}
}

type RunProgramRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *WiredSyncHandler) RunProgram(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get username and password
	var req RunProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	// Run the command with username and password as arguments
	cmd := exec.Command("go", "run", "wiredsync/main.go", req.Username, req.Password)

	// Capture output and handle errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running program: %v\nOutput: %s", err, string(output)), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Program executed successfully", "output": "%s"}`, string(output))))
}
