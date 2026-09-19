package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	OpenSearch OpenSearchConfig
	Redis      RedisConfig
	Tracing    TracingConfig
}

type ServerConfig struct {
	Port string
}

type OpenSearchConfig struct {
	URL      string
	Username string
	Password string
}

type RedisConfig struct {
	URL string
}

type TracingConfig struct {
	ServiceName string
	Endpoint    string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		OpenSearch: OpenSearchConfig{
			URL:      os.Getenv("OPENSEARCH_URL"),
			Username: os.Getenv("OPENSEARCH_USERNAME"),
			Password: os.Getenv("OPENSEARCH_PASSWORD"),
		},
		Redis: RedisConfig{
			URL: os.Getenv("REDIS_URL"),
		},
		Tracing: TracingConfig{
			ServiceName: os.Getenv("OTEL_SERVICE_NAME"),
			Endpoint:    os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		},
	}
}
