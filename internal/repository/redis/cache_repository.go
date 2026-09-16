package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"library-app-search/internal/domain"
)

type CacheRepository struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewRedisCacheRepository(client *goredis.Client, ttl time.Duration) *CacheRepository {
	return &CacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func buildCacheKey(params domain.SearchParams) string {
	return fmt.Sprintf(
		"search:chapters:%s:page=%d:size=%d",
		params.Query,
		params.Page,
		params.Size,
	)
}

func (r *CacheRepository) FetchChaptersCache(params domain.SearchParams) (domain.SearchResult, bool, error) {
	key := buildCacheKey(params)

	data, err := r.client.Get(context.Background(), key).Result()

	if err != nil {
		if err == goredis.Nil {
			return domain.SearchResult{}, false, nil
		}

		return domain.SearchResult{}, false, err
	}

	var result domain.SearchResult

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return domain.SearchResult{}, false, err
	}

	return result, true, nil
}

func (r *CacheRepository) PutChaptersCache(params domain.SearchParams, result domain.SearchResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	key := buildCacheKey(params)
	err = r.client.Set(
		context.Background(),
		key,
		data,
		r.ttl,
	).Err()

	if err != nil {
		return err
	}

	return nil
}
