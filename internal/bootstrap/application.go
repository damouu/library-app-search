package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"library-app-search/internal/application"
	"library-app-search/internal/config"
	elasticsearchClient "library-app-search/internal/elasticsearch"
	"library-app-search/internal/handler"
	"library-app-search/internal/health"
	elasticsearchRepository "library-app-search/internal/repository/elasticsearch"
	redisRepository "library-app-search/internal/repository/redis"
	"library-app-search/internal/router"

	goredis "github.com/redis/go-redis/v9"
)

type Application struct {
	Router      http.Handler
	RedisClient *goredis.Client
	Port        string
}

func New(cfg config.Config) (*Application, error) {
	elasticsearch, err := elasticsearchClient.CreateClient(&cfg.Elasticsearch)
	if err != nil {
		return nil, fmt.Errorf("create Elasticsearch client: %w", err)
	}

	searchRepository := elasticsearchRepository.NewSearchRepository(elasticsearch)

	redisClient, err := redisRepository.CreateClient(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("create Redis client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = redisClient.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	cacheRepository := redisRepository.NewRedisCacheRepository(
		redisClient,
		5*time.Minute,
	)

	searchService := application.NewSearchService(
		cacheRepository,
		searchRepository,
	)

	searchHandler := handler.NewSearchHandler(searchService)

	healthChecker := elasticsearchClient.NewHealthChecker(elasticsearch)
	healthHandler := health.NewHandler(healthChecker)

	httpRouter := router.NewRouter(
		healthHandler,
		searchHandler,
	)

	return &Application{
		Router:      httpRouter,
		RedisClient: redisClient,
		Port:        cfg.Server.Port,
	}, nil
}
