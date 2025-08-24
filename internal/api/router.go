package api

import (
	"net/http"

	"github.com/keshupandre/img-to-pdf-backend/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/api/health", handlers.HealthHandler)
	r.Post("/api/convert", handlers.ConvertHandler)

	return r
}
