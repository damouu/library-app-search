package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"library-app-search/internal/application"
	"library-app-search/internal/config"
	"library-app-search/internal/handler"
	"library-app-search/internal/health"
	opensearchClient "library-app-search/internal/opensearch"
	opensearchRepository "library-app-search/internal/repository/opensearch"
	redisRepository "library-app-search/internal/repository/redis"
	"library-app-search/internal/router"
)

// Application contains the fully initialized dependencies required by the HTTP server.
type Application struct {
	Router      http.Handler
	RedisClient *redis.Client
	Port        string
}

// New builds the application dependency graph from the provided configuration.
func New(cfg config.Config) (*Application, error) {
	openSearch, err := opensearchClient.CreateClient(&cfg.OpenSearch)
	if err != nil {
		return nil, fmt.Errorf("create OpenSearch client: %w", err)
	}

	searchRepository := opensearchRepository.NewSearchRepository(openSearch)

	redisClient, err := redisRepository.CreateClient(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("create Redis client: %w", err)
	}

	// Fail fast if Redis is unavailable during application startup.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = redisClient.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	// Search results are cached for five minutes to reduce repeated OpenSearch queries.
	cacheRepository := redisRepository.NewRedisCacheRepository(redisClient, 5*time.Minute)

	searchService := application.NewSearchService(cacheRepository, searchRepository)

	searchHandler := handler.NewSearchHandler(searchService)

	healthChecker := opensearchClient.NewHealthChecker(openSearch)
	healthHandler := health.NewHandler(healthChecker)

	httpRouter := router.NewRouter(healthHandler, searchHandler)

	return &Application{
		Router:      httpRouter,
		RedisClient: redisClient,
		Port:        cfg.Server.Port,
	}, nil
}
