package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"img-to-pdf-converter/internal/config"
	"img-to-pdf-converter/internal/models"
	"img-to-pdf-converter/internal/services"
)

// Handler holds all the dependencies for HTTP handlers
type Handler struct {
	config      *config.Config
	pdfService  *services.PDFService
	fileService *services.FileService
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, pdfService *services.PDFService, fileService *services.FileService) *Handler {
	return &Handler{
		config:      cfg,
		pdfService:  pdfService,
		fileService: fileService,
	}
}

// sendErrorResponse sends an error response in JSON format
func (h *Handler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helper function to get map keys for logging
func getStringMapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to get the first non-empty string or return default
func getFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// Helper function to get file map keys for logging
func getFileMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
