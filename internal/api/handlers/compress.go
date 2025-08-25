package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/keshupandre/img-to-pdf-backend/internal/services"

	"github.com/google/uuid"
)

func CompressHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // 20 MB
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["pdf"]
	if len(files) == 0 {
		http.Error(w, "No PDF uploaded", http.StatusBadRequest)
		return
	}

	// Quality param (default 75)
	quality := 75
	if qStr := r.URL.Query().Get("quality"); qStr != "" {
		if q, err := strconv.Atoi(qStr); err == nil && q >= 1 && q <= 100 {
			quality = q
		}
	}

	// Ensure uploads folder exists
	os.MkdirAll("uploads", os.ModePerm)

	fileHeader := files[0]
	file, _ := fileHeader.Open()
	defer file.Close()

	tmpPath := filepath.Join("uploads", fileHeader.Filename)
	out, _ := os.Create(tmpPath)
	defer out.Close()
	_, _ = io.Copy(out, file)

	// Generate unique compressed file
	uniqueID := uuid.New().String()
	outputFilename := fmt.Sprintf("%s-compressed.pdf", uniqueID)
	outputPath := filepath.Join("uploads", outputFilename)

	err = services.CompressPDF(tmpPath, outputPath, quality)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Failed to compress PDF", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFilename))

	pdfFile, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, "Could not open compressed PDF", http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	io.Copy(w, pdfFile)
}
