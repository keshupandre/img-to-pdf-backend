package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadHandler handles PDF file downloads
// This is a new implementation that will replace the one in handlers.go
func (h *Handler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Download Handler  Called ===")

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		h.sendErrorResponse(w, "No filename specified", http.StatusBadRequest)
		return
	}

	// Security check: prevent directory traversal
	if filepath.Dir(filename) != "." {
		h.sendErrorResponse(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.config.PDF.OutputDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.sendErrorResponse(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Serve the file
	http.ServeFile(w, r, filePath)
	log.Printf("File served: %s", filePath)
}
