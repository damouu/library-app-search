package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	Elasticsearch ElasticsearchConfig
	Redis         RedisConfig
	Tracing       TracingConfig
}

type ServerConfig struct {
	Port string
}

type ElasticsearchConfig struct {
	URL    string
	APIKey string
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
		Server:        ServerConfig{Port: os.Getenv("SERVER_PORT")},
		Elasticsearch: ElasticsearchConfig{URL: os.Getenv("ELASTICSEARCH_URL"), APIKey: os.Getenv("ELASTICSEARCH_API_KEY")},
		Redis:         RedisConfig{URL: os.Getenv("REDIS_URL")},
		Tracing:       TracingConfig{ServiceName: os.Getenv("OTEL_SERVICE_NAME"), Endpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")},
	}
}
