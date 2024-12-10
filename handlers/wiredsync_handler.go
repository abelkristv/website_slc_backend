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
	var req RunProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("go", "run", "wiredsync/main.go", req.Username, req.Password)

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

	// Start the process without blocking
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start program: %v", err), http.StatusInternalServerError)
		return
	}

	// Use goroutines to process stdout and stderr
	w.Header().Set("Content-Type", "text/plain")
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Fprintf(w, "[STDOUT] %s\n", scanner.Text())
			w.(http.Flusher).Flush()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(w, "[STDERR] %s\n", scanner.Text())
			w.(http.Flusher).Flush()
		}
	}()

	// Wait for the process to complete
	if err := cmd.Wait(); err != nil {
		http.Error(w, fmt.Sprintf("Program finished with an error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Program execution completed successfully\n")
}
