package main

import (
	"library-app-search/internal/config"
	elasticsearch2 "library-app-search/internal/elasticsearch"
	"library-app-search/internal/health"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Load()
	client, err := elasticsearch2.CreateClient(&cfg.Elasticsearch)

	if err != nil {
		log.Fatal(err)
	}

	healthChecker := elasticsearch2.NewHealthChecker(client)

	healthHandler := health.NewHandler(healthChecker)

	router := chi.NewRouter()

	router.Get("/health/live", health.LiveHandler)

	router.Get("/health/ready", healthHandler.ReadyHandler)

	log.Println("Search API listening on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
