package main

import (
	"library-app-search/internal/application"
	"library-app-search/internal/config"
	elasticsearch2 "library-app-search/internal/elasticsearch"
	"library-app-search/internal/handler"
	"library-app-search/internal/health"
	elasticsearchRepository "library-app-search/internal/repository/elasticsearch"
	redisRepository "library-app-search/internal/repository/redis"
	"library-app-search/internal/router"
	"log"
	"net/http"
)

func main() {

	cfg := config.Load()
	client, err := elasticsearch2.CreateClient(&cfg.Elasticsearch)

	if err != nil {
		log.Fatal(err)
	}

	searchRepository := elasticsearchRepository.NewSearchRepository(client)

	cacheRepository := redisRepository.NewNoopCacheRepository()

	searchService := application.NewSearchService(
		cacheRepository,
		searchRepository,
	)

	searchHandler := handler.NewSearchHandler(searchService)

	healthChecker := elasticsearch2.NewHealthChecker(client)

	healthHandler := health.NewHandler(healthChecker)

	httpRouter := router.NewRouter(healthHandler, searchHandler)

	log.Println("Search API listening on :8080")

	if err := http.ListenAndServe(":8080", httpRouter); err != nil {
		log.Fatal(err)
	}
}
