package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"library-app-search/internal/handler"
	"library-app-search/internal/health"
)

func NewRouter(h *health.Handler, searchHandler *handler.SearchHandler) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health/live", health.LiveHandler)

	router.Get("/health/ready", h.ReadyHandler)

	router.Get("/search", searchHandler.Search)

	return router

}
