package redis

import (
	goredis "github.com/redis/go-redis/v9"

	"library-app-search/internal/config"
)

func CreateClient(cfg *config.RedisConfig) (*goredis.Client, error) {
	options, err := goredis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	client := goredis.NewClient(options)

	return client, nil
}
