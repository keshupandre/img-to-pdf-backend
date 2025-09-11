package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"time"
)

// GenerateTimestampedFilename generates a filename with timestamp
func GenerateTimestampedFilename(prefix, extension string) string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s.%s", prefix, timestamp, extension)
}

// GetImageDimensions returns the dimensions of an image file
func GetImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

// SanitizeFilename removes dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Keep only the base filename, remove path
	base := filepath.Base(filename)

	// In a real implementation, you might want to remove/replace
	// dangerous characters, but for now we'll just return the base
	return base
}

// CreateTempDir creates a temporary directory with a unique name
func CreateTempDir(baseDir, prefix string) (string, error) {
	timestamp := time.Now().UnixNano()
	tempDirName := fmt.Sprintf("%s_%d", prefix, timestamp)
	tempDirPath := filepath.Join(baseDir, tempDirName)

	err := os.MkdirAll(tempDirPath, 0755)
	if err != nil {
		return "", err
	}

	return tempDirPath, nil
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filepath string) (int64, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
