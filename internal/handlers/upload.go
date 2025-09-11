package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"img-to-pdf-converter/internal/models"
	"img-to-pdf-converter/internal/services"
)

// UploadHandler handles file uploads and PDF conversion
// This is a new implementation that will replace the one in handlers.go
func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Upload Handler  Called ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
	log.Printf("Content-Length: %s", r.Header.Get("Content-Length"))

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		h.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(h.config.Upload.MaxFileSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		h.sendErrorResponse(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	log.Printf("Multipart form parsed successfully")
	log.Printf("Form keys: %v", getStringMapKeys(r.MultipartForm.Value))
	log.Printf("File keys: %v", getFileMapKeys(r.MultipartForm.File))

	// Parse conversion options from query parameters or form data
	positionValue := getFirstNonEmpty(r.URL.Query().Get("position"), r.FormValue("position"))
	if positionValue == "" {
		positionValue = "center"
	}

	orientationValue := getFirstNonEmpty(r.URL.Query().Get("orientation"), r.FormValue("orientation"))
	if orientationValue == "" {
		orientationValue = "P"
	}

	options := services.ConversionOptions{
		Fit:         r.URL.Query().Get("fit") == "true" || r.FormValue("fit") == "true",
		Position:    positionValue,
		Orientation: orientationValue,
	}

	log.Printf("Conversion options: fit=%t, position=%s, orientation=%s", options.Fit, options.Position, options.Orientation)

	// Get uploaded files - try both 'images' and 'files' field names
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		files = r.MultipartForm.File["files"]
	}
	if len(files) == 0 {
		log.Printf("No files found in form data")
		h.sendErrorResponse(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	log.Printf("Found %d files", len(files))
	for i, file := range files {
		log.Printf("File %d: name=%s, size=%d, header=%v", i, file.Filename, file.Size, file.Header)
	}

	// Validate files
	if err := h.fileService.ValidateFiles(files); err != nil {
		log.Printf("File validation failed: %v", err)
		h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert images to PDF with options
	pdfPath, err := h.pdfService.ConvertImagesToPDFWithOptions(files, options)
	if err != nil {
		log.Printf("PDF conversion failed: %v", err)
		h.sendErrorResponse(w, "Failed to convert images to PDF", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := models.UploadResponse{
		Success: true,
		PDFFile: filepath.Base(pdfPath),
		Message: "Images converted to PDF successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Printf("Upload completed successfully, PDF: %s", pdfPath)
}
