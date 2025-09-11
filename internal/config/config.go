package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	CORS   CORSConfig
	Upload UploadConfig
	PDF    PDFConfig
	App    AppConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port  string
	Host  string
	Debug bool
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// UploadConfig holds upload-related configuration
type UploadConfig struct {
	MaxFileSize  int64
	MaxFiles     int
	AllowedTypes []string
	TempDir      string
	UploadDir    string
}

// PDFConfig holds PDF generation configuration
type PDFConfig struct {
	OutputDir   string
	PageFormat  string
	Orientation string
	Unit        string
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

// Load creates and returns a new Config instance with values from environment variables or defaults
func Load() *Config {
	godotenv.Load()
	return &Config{
		Server: ServerConfig{
			Port:  getEnvOrDefault("PORT", "8080"),
			Host:  getEnvOrDefault("HOST", "localhost"),
			Debug: getEnvBoolOrDefault("DEBUG", true),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
				"http://127.0.0.1:3000",
				"http://localhost:3001",
				"http://127.0.0.1:3001",
			},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		},
		Upload: UploadConfig{
			MaxFileSize:  getEnvIntOrDefault("MAX_FILE_SIZE", 10*1024*1024), // 10MB
			MaxFiles:     int(getEnvIntOrDefault("MAX_FILES", 10)),
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp"},
			TempDir:      getEnvOrDefault("TEMP_DIR", "./temp"),
			UploadDir:    getEnvOrDefault("UPLOAD_DIR", "./uploads"),
		},
		PDF: PDFConfig{
			OutputDir:   getEnvOrDefault("PDF_OUTPUT_DIR", "./output"),
			PageFormat:  getEnvOrDefault("PDF_PAGE_FORMAT", "A4"),
			Orientation: getEnvOrDefault("PDF_ORIENTATION", "P"),
			Unit:        getEnvOrDefault("PDF_UNIT", "mm"),
		},
		App: AppConfig{
			Name:        getEnvOrDefault("APP_NAME", "Image to PDF Converter"),
			Version:     getEnvOrDefault("APP_VERSION", "1.0.0"),
			Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		},
	}
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
