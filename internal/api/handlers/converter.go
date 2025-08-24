package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/keshupandre/img-to-pdf-backend/internal/services"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {

	cwd, _ := os.Getwd()

	uploadsDir := filepath.Join(cwd, "uploads")
	os.MkdirAll(uploadsDir, os.ModePerm)

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["images"] // form field name = "images"
	if len(files) == 0 {
		http.Error(w, "No images uploaded", http.StatusBadRequest)
		return
	}

	var imagePaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusInternalServerError)
			return
		}

		tmpPath := filepath.Join(uploadsDir, fileHeader.Filename)
		out, err := os.Create(tmpPath)
		if err != nil {
			file.Close()
			http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
			return
		}

		_, err = out.ReadFrom(file)
		file.Close()
		out.Close()
		if err != nil {
			http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
			return
		}
		imagePaths = append(imagePaths, tmpPath)
	}

	// Convert to PDF
	outputPath := filepath.Join(uploadsDir, "output.pdf")
	err = services.ImagesToPDF(imagePaths, outputPath)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Response (URL path should be relative to /uploads/)
	response := map[string]string{"pdf_url": "/uploads/output.pdf"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
