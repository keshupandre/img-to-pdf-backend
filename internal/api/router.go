package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/keshupandre/img-to-pdf-backend/internal/api/handlers"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/api/health", handlers.HealthHandler)
	r.Post("/api/convert", handlers.ConvertHandler)
	r.Post("/api/compress", handlers.CompressHandler)

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("../../uploads"))))

	return r
}
