package handlers

import (
	"encoding/json"
	"net/http"

	"img-to-pdf-converter/internal/models"
)

// HealthHandler handles health check requests
// This is a new implementation that will replace the one in handlers.go
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:  "healthy",
		Message: "Image to PDF Converter is running",
		Version: h.config.App.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
