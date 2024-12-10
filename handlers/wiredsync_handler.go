package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	var req RunProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	// Prepare the command
	cmd := exec.Command("go", "run", "wiredsync/main.go", req.Username, req.Password)

	// Bind stdout and stderr to the OS terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start program: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the response immediately
	fmt.Fprint(w, "Program is running. Check terminal logs for output.\n")

	// Wait for the process to complete
	if err := cmd.Wait(); err != nil {
		// Log the error to the terminal but do not write a second response
		fmt.Fprintf(os.Stderr, "Program finished with an error: %v\n", err)
		return
	}

	// Log successful completion to the terminal
	fmt.Fprintln(os.Stdout, "Program execution completed successfully")
}
