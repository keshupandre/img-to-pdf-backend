package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"img-to-pdf-converter/internal/config"
)

// FileService handles file operations
type FileService struct {
	config *config.Config
}

// NewFileService creates a new file service instance
func NewFileService(cfg *config.Config) *FileService {
	return &FileService{
		config: cfg,
	}
}

// ValidateFile validates an uploaded file
func (s *FileService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.config.Upload.MaxFileSize {
		return fmt.Errorf("file %s is too large: %d bytes (max: %d bytes)",
			fileHeader.Filename, fileHeader.Size, s.config.Upload.MaxFileSize)
	}

	// Check file type
	contentType := fileHeader.Header.Get("Content-Type")
	if !s.isAllowedType(contentType) {
		return fmt.Errorf("file %s has unsupported type: %s", fileHeader.Filename, contentType)
	}

	return nil
}

// ValidateFiles validates multiple uploaded files
func (s *FileService) ValidateFiles(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return fmt.Errorf("no files provided")
	}

	if len(files) > s.config.Upload.MaxFiles {
		return fmt.Errorf("too many files: %d (max: %d)", len(files), s.config.Upload.MaxFiles)
	}

	for _, file := range files {
		if err := s.ValidateFile(file); err != nil {
			return err
		}
	}

	return nil
}

// isAllowedType checks if the content type is allowed
func (s *FileService) isAllowedType(contentType string) bool {
	for _, allowedType := range s.config.Upload.AllowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

// GetFileExtension returns the file extension from filename
func (s *FileService) GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsImageFile checks if a file is an image based on its extension
func (s *FileService) IsImageFile(filename string) bool {
	ext := s.GetFileExtension(filename)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// EnsureDirectoryExists creates a directory if it doesn't exist
func (s *FileService) EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// CleanupFile removes a file
func (s *FileService) CleanupFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}

// CleanupDirectory removes a directory and all its contents
func (s *FileService) CleanupDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}
