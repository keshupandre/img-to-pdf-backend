package models

import "mime/multipart"

// UploadRequest represents the incoming file upload request
type UploadRequest struct {
	Files []*multipart.FileHeader `json:"files"`
}

// UploadResponse represents the response after successful upload and conversion
type UploadResponse struct {
	Success bool   `json:"success"`
	PDFFile string `json:"pdfFile"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}

// ImageFile represents a processed image file
type ImageFile struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	TempPath string `json:"temp_path,omitempty"`
}

// ConversionJob represents a PDF conversion job
type ConversionJob struct {
	ID        string      `json:"id"`
	Images    []ImageFile `json:"images"`
	PDFPath   string      `json:"pdf_path"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
}
