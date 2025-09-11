package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

	"img-to-pdf-converter/internal/config"
)

// PDFService handles PDF conversion operations
type PDFService struct {
	config *config.Config
}

// ConversionOptions holds the conversion parameters
type ConversionOptions struct {
	Fit         bool   `json:"fit"`         // Fit small images to page
	Position    string `json:"position"`    // Image positioning
	Orientation string `json:"orientation"` // PDF orientation
}

// Margins defines page margins
type Margins struct {
	Top, Right, Bottom, Left float64
}

// NewPDFService creates a new PDF service instance
func NewPDFService(cfg *config.Config) *PDFService {
	return &PDFService{
		config: cfg,
	}
}

// ConvertImagesToPDF converts uploaded images to a single PDF file with options
func (s *PDFService) ConvertImagesToPDF(files []*multipart.FileHeader) (string, error) {
	return s.ConvertImagesToPDFWithOptions(files, ConversionOptions{
		Fit:         false,
		Position:    "center",
		Orientation: "P", // Portrait by default
	})
}

// ConvertImagesToPDFWithOptions converts uploaded images to a single PDF file with conversion options
func (s *PDFService) ConvertImagesToPDFWithOptions(files []*multipart.FileHeader, options ConversionOptions) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no files provided")
	}

	// Create output directory if it doesn't exist
	if err := s.ensureDirectory(s.config.PDF.OutputDir); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create temporary directory for processing
	tempDir := filepath.Join(s.config.Upload.TempDir, fmt.Sprintf("conversion_%d", time.Now().Unix()))
	if err := s.ensureDirectory(tempDir); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer s.cleanupDirectory(tempDir)

	// Save uploaded files to temp directory
	var imagePaths []string
	for i, fileHeader := range files {
		tempPath, err := s.saveUploadedFile(fileHeader, tempDir, fmt.Sprintf("image_%d_%s", i, fileHeader.Filename))
		if err != nil {
			return "", fmt.Errorf("failed to save file %s: %v", fileHeader.Filename, err)
		}
		imagePaths = append(imagePaths, tempPath)
	}

	// Generate PDF with options
	pdfPath, err := s.generatePDFWithOptions(imagePaths, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %v", err)
	}

	return pdfPath, nil
}

// saveUploadedFile saves an uploaded file to the specified directory
func (s *PDFService) saveUploadedFile(fileHeader *multipart.FileHeader, destDir, filename string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	destPath := filepath.Join(destDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	log.Printf("Saved file: %s (size: %d bytes)", destPath, fileHeader.Size)
	return destPath, nil
}

// generatePDFWithOptions creates a PDF from the provided image paths with conversion options
func (s *PDFService) generatePDFWithOptions(imagePaths []string, options ConversionOptions) (string, error) {
	if len(imagePaths) == 0 {
		return "", fmt.Errorf("no images provided")
	}

	// Validate and set default orientation
	orientation := options.Orientation
	if orientation != "L" && orientation != "P" {
		orientation = "P" // Default to Portrait
	}

	// Create PDF with specified orientation
	pdf := gofpdf.New(orientation, "mm", "A4", "")

	// Calculate page dimensions based on orientation
	var pageW, pageH float64
	if orientation == "L" {
		pageW, pageH = 297.0, 210.0 // Landscape A4
	} else {
		pageW, pageH = 210.0, 297.0 // Portrait A4
	}

	// Define margins
	margins := Margins{
		Top:    10,
		Right:  10,
		Bottom: 10,
		Left:   10,
	}

	// Calculate usable area
	usableW := pageW - (margins.Left + margins.Right)
	usableH := pageH - (margins.Top + margins.Bottom)

	for i, imagePath := range imagePaths {
		log.Printf("Processing image %d: %s", i+1, imagePath)

		// Check if file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			log.Printf("Warning: Image file does not exist: %s", imagePath)
			continue
		}

		// Add new page
		pdf.AddPage()

		// Register image and get dimensions
		info := pdf.RegisterImage(imagePath, "")
		if info == nil {
			log.Printf("Warning: Failed to register image: %s", imagePath)
			continue
		}

		imgW, imgH := info.Extent()
		log.Printf("Original image dimensions: %.2f x %.2f", imgW, imgH)

		// Calculate new dimensions
		newW, newH := s.calculateOptimalDimensions(imgW, imgH, usableW, usableH, options.Fit)
		log.Printf("Calculated dimensions: %.2f x %.2f", newW, newH)

		// Calculate position
		x, y := s.calculatePosition(options.Position, usableW, usableH, newW, newH)
		x += margins.Left
		y += margins.Top

		log.Printf("Image position: %.2f, %.2f", x, y)

		// Add image to PDF
		pdf.ImageOptions(imagePath, x, y, newW, newH, false, gofpdf.ImageOptions{}, 0, "")
		log.Printf("Added image to PDF: %s", imagePath)
	}

	// Generate output filename
	timestamp := time.Now().Format("20060102_150405")
	outputFilename := fmt.Sprintf("converted_images_%s.pdf", timestamp)
	outputPath := filepath.Join(s.config.PDF.OutputDir, outputFilename)

	// Save PDF
	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %v", err)
	}

	log.Printf("PDF saved: %s", outputPath)
	return outputPath, nil
}

// calculateOptimalDimensions calculates optimal image dimensions based on fit option
func (s *PDFService) calculateOptimalDimensions(imgW, imgH, usableW, usableH float64, fit bool) (float64, float64) {
	newW, newH := imgW, imgH

	// If image is larger than usable area, scale it down
	if imgW > usableW || imgH > usableH {
		ratio := s.min(usableW/imgW, usableH/imgH)
		newW = imgW * ratio
		newH = imgH * ratio
	} else if fit {
		// If fit is enabled and image is smaller, scale it up to fit the area
		ratio := s.min(usableW/imgW, usableH/imgH)
		newW = imgW * ratio
		newH = imgH * ratio
	}

	return newW, newH
}

// calculatePosition calculates image position based on position option
func (s *PDFService) calculatePosition(pos string, areaW, areaH, imgW, imgH float64) (float64, float64) {
	var x, y float64

	switch strings.ToLower(pos) {
	case "top-left":
		x, y = 0, 0
	case "top-center", "top":
		x, y = (areaW-imgW)/2, 0
	case "top-right":
		x, y = areaW-imgW, 0
	case "center-left", "left":
		x, y = 0, (areaH-imgH)/2
	case "center", "center-center", "":
		x, y = (areaW-imgW)/2, (areaH-imgH)/2
	case "center-right", "right":
		x, y = areaW-imgW, (areaH-imgH)/2
	case "bottom-left":
		x, y = 0, areaH-imgH
	case "bottom-center", "bottom":
		x, y = (areaW-imgW)/2, areaH-imgH
	case "bottom-right":
		x, y = areaW-imgW, areaH-imgH
	default:
		// Default to center
		x, y = (areaW-imgW)/2, (areaH-imgH)/2
	}

	return x, y
}

// min returns the minimum of two float64 values
func (s *PDFService) min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ensureDirectory creates a directory if it doesn't exist
func (s *PDFService) ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// cleanupDirectory removes a directory and its contents
func (s *PDFService) cleanupDirectory(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("Warning: Failed to cleanup directory %s: %v", dir, err)
	}
}
