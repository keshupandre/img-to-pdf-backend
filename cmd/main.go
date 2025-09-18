package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"img-to-pdf-converter/internal/config"
	"img-to-pdf-converter/internal/handlers"
	"img-to-pdf-converter/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting %s v%s in %s mode at %v PORT", cfg.App.Name, cfg.App.Version, cfg.App.Environment, cfg.Server.Port)

	// Initialize services
	fileService := services.NewFileService(cfg)
	pdfService := services.NewPDFService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(cfg, pdfService, fileService)

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Define routes
	router.Post("/upload", handler.UploadHandler)
	router.Get("/download", handler.DownloadHandler)
	router.Get("/health", handler.HealthHandler)

	// Add a simple root endpoint
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Image to PDF Converter API", "version": "` + cfg.App.Version + `"}`))
	})

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
		Debug:            cfg.Server.Debug,
	})

	// Create server with CORS middleware
	handler_with_cors := c.Handler(router)

	// Ensure required directories exist
	if err := fileService.EnsureDirectoryExists(cfg.Upload.TempDir); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	if err := fileService.EnsureDirectoryExists(cfg.Upload.UploadDir); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := fileService.EnsureDirectoryExists(cfg.PDF.OutputDir); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on http://%s", serverAddr)
	log.Printf("CORS allowed origins: %v", cfg.CORS.AllowedOrigins)
	log.Printf("Upload config: MaxFileSize=%d bytes, MaxFiles=%d", cfg.Upload.MaxFileSize, cfg.Upload.MaxFiles)
	log.Printf("Allowed file types: %v", cfg.Upload.AllowedTypes)

	if err := http.ListenAndServe(serverAddr, handler_with_cors); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
