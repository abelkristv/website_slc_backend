package handlers

import (
	"bufio"
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

	// Create pipes to capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to capture stdout: %v", err), http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to capture stderr: %v", err), http.StatusInternalServerError)
		return
	}

	// Start the command execution
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start program: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Stream stdout to the client
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Fprintf(w, "[STDOUT] %s\n", scanner.Text())
			w.(http.Flusher).Flush() // Ensure the data is sent immediately
		}
	}()

	// Stream stderr to the client
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(w, "[STDERR] %s\n", scanner.Text())
			w.(http.Flusher).Flush() // Ensure the data is sent immediately
		}
	}()

	// Wait for the program to finish
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(w, "Program finished with an error: %v\n", err)
		w.(http.Flusher).Flush()
		return
	}

	// Notify the client that the program has completed
	fmt.Fprint(w, "Program execution completed successfully\n")
	w.(http.Flusher).Flush()
}
