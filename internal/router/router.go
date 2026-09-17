package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"library-app-search/internal/handler"
	"library-app-search/internal/health"
)

// NewRouter builds the HTTP router and registers the service endpoints.
func NewRouter(
	h *health.Handler,
	searchHandler *handler.SearchHandler,
) chi.Router {
	router := chi.NewRouter()

	// Add a unique ID to each request for log correlation.
	router.Use(middleware.RequestID)

	// Log incoming HTTP requests and their response status.
	router.Use(middleware.Logger)

	// Recover from panics and prevent them from crashing the HTTP server.
	router.Use(middleware.Recoverer)

	router.Route("/health", func(r chi.Router) {
		// Liveness check: verifies that the service process is running.
		r.Get("/live", health.LiveHandler)

		// Readiness check: verifies that the required dependencies are available.
		r.Get("/ready", h.ReadyHandler)
	})

	// Search chapters using the Search API.
	router.Get("/search", searchHandler.Search)

	return router
}
