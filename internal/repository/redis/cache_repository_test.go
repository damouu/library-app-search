package redis

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	miniRedis "github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"

	"library-app-search/internal/domain"
)

func newTestCacheRepository(t *testing.T, ttl time.Duration) (*CacheRepository, *miniRedis.Miniredis, *goredis.Client) {
	t.Helper()

	server, err := miniRedis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	client := goredis.NewClient(&goredis.Options{
		Addr: server.Addr(),
	})

	repository := NewRedisCacheRepository(client, ttl)

	t.Cleanup(func() {
		_ = client.Close()
		server.Close()
	})

	return repository, server, client
}

func TestCacheRepository_FetchChaptersCache_ReturnsMissWhenKeyDoesNotExist(t *testing.T) {
	repository, _, _ := newTestCacheRepository(t, 5*time.Minute)

	params := domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	}

	result, found, err := repository.FetchChaptersCache(params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found {
		t.Fatal("expected cache miss, got cache hit")
	}

	if !reflect.DeepEqual(result, domain.SearchResult{}) {
		t.Fatalf("expected empty result, got %+v", result)
	}
}

func TestCacheRepository_PutAndFetchChaptersCache(t *testing.T) {
	repository, _, _ := newTestCacheRepository(t, 5*time.Minute)

	params := domain.SearchParams{
		Query: "one piece",
		Page:  2,
		Size:  3,
	}

	expectedResult := domain.SearchResult{
		Items: []domain.Chapter{
			{
				ChapterUUID: "chapter-1",
				Title:       "One Piece",
			},
		},
		Page:  2,
		Size:  3,
		Total: 1,
	}

	err := repository.PutChaptersCache(params, expectedResult)
	if err != nil {
		t.Fatalf("expected no error while writing cache, got %v", err)
	}

	result, found, err := repository.FetchChaptersCache(params)
	if err != nil {
		t.Fatalf("expected no error while fetching cache, got %v", err)
	}

	if !found {
		t.Fatal("expected cache hit, got cache miss")
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf(
			"unexpected cached result:\ngot:  %+v\nwant: %+v",
			result,
			expectedResult,
		)
	}
}

func TestCacheRepository_PutChaptersCache_UsesExpectedCacheKey(t *testing.T) {
	repository, _, client := newTestCacheRepository(t, 5*time.Minute)

	params := domain.SearchParams{
		Query: "ワンピース",
		Page:  1,
		Size:  3,
	}

	expectedResult := domain.SearchResult{
		Items: []domain.Chapter{},
		Page:  1,
		Size:  3,
		Total: 0,
	}

	err := repository.PutChaptersCache(params, expectedResult)
	if err != nil {
		t.Fatalf("expected no error while writing cache, got %v", err)
	}

	expectedKey := "search:chapters:ワンピース:page=1:size=3"

	if exists, err := client.Exists(context.Background(), expectedKey).Result(); err != nil {
		t.Fatalf("failed to check cache key: %v", err)
	} else if exists != 1 {
		t.Fatalf("expected cache key %q to exist", expectedKey)
	}
}

func TestCacheRepository_PutChaptersCache_ExpiresAfterTTL(t *testing.T) {
	ttl := 5 * time.Minute

	repository, miniRedis, _ := newTestCacheRepository(t, ttl)

	params := domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	}

	expectedResult := domain.SearchResult{
		Items: []domain.Chapter{
			{
				ChapterUUID: "chapter-1",
				Title:       "One Piece",
			},
		},
		Page:  1,
		Size:  10,
		Total: 1,
	}

	err := repository.PutChaptersCache(params, expectedResult)
	if err != nil {
		t.Fatalf("expected no error while writing cache, got %v", err)
	}

	_, found, err := repository.FetchChaptersCache(params)
	if err != nil {
		t.Fatalf("expected no error while fetching cache, got %v", err)
	}

	if !found {
		t.Fatal("expected cache hit before TTL expiration")
	}

	miniRedis.FastForward(ttl + time.Second)

	_, found, err = repository.FetchChaptersCache(params)
	if err != nil {
		t.Fatalf("expected no error after TTL expiration, got %v", err)
	}

	if found {
		t.Fatal("expected cache miss after TTL expiration")
	}
}

func TestCacheRepository_FetchChaptersCache_ReturnsErrorWhenRedisFails(t *testing.T) {
	repository, miniRedis, _ := newTestCacheRepository(t, 5*time.Minute)

	miniRedis.Close()

	params := domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	}

	_, found, err := repository.FetchChaptersCache(params)

	if err == nil {
		t.Fatal("expected Redis error, got nil")
	}

	if found {
		t.Fatal("expected found to be false when Redis fails")
	}
}

func TestCacheRepository_PutChaptersCache_ReturnsErrorWhenRedisFails(t *testing.T) {
	repository, miniRedis, _ := newTestCacheRepository(t, 5*time.Minute)

	miniRedis.Close()

	params := domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	}

	result := domain.SearchResult{
		Items: []domain.Chapter{},
		Page:  1,
		Size:  10,
		Total: 0,
	}

	err := repository.PutChaptersCache(params, result)

	if err == nil {
		t.Fatal("expected Redis error, got nil")
	}
}

func TestCacheRepository_FetchChaptersCache_ReturnsErrorWhenCachedDataIsInvalid(t *testing.T) {
	repository, _, client := newTestCacheRepository(t, 5*time.Minute)

	params := domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	}

	key := buildCacheKey(params)

	err := client.Set(
		context.Background(),
		key,
		"invalid-json",
		5*time.Minute,
	).Err()
	if err != nil {
		t.Fatalf("failed to prepare invalid cache data: %v", err)
	}

	_, found, err := repository.FetchChaptersCache(params)

	if err == nil {
		t.Fatal("expected JSON unmarshal error, got nil")
	}

	if found {
		t.Fatal("expected found to be false when cached data is invalid")
	}

	if errors.Is(err, goredis.Nil) {
		t.Fatalf("expected JSON error, got Redis nil error: %v", err)
	}
}
