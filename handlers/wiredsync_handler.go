package handlers

import (
	"bufio"
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

	// Set the response headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Flush the headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create pipes for redirecting os.Stdout and os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	_, wStderr, _ := os.Pipe()

	// Replace os.Stdout and os.Stderr with the writers
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	os.Stdout = wStdout
	os.Stderr = wStderr

	// Restore os.Stdout and os.Stderr after the function exits
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	// Start goroutines to send logs to SSE
	go func() {
		scanner := bufio.NewScanner(rStdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Fprintf(w, "data: [STDOUT] %s\n\n", line) // Send to SSE
			flusher.Flush()
			fmt.Fprintln(originalStdout, line) // Print to terminal
		}
	}()

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
