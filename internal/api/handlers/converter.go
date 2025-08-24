package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

	uniqueID := uuid.New().String()
	outputFilename := fmt.Sprintf("%s.pdf", uniqueID)
	outputPath := filepath.Join("uploads", outputFilename)

	err = services.ImagesToPDF(imagePaths, outputPath)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")

	pdfFile, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, "Could not open generated PDF", http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	io.Copy(w, pdfFile)
}
