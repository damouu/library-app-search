package redis

import "library-app-search/internal/domain"

type NoopCacheRepository struct{}

func NewNoopCacheRepository() *NoopCacheRepository {
	return &NoopCacheRepository{}
}

func (c *NoopCacheRepository) FetchChaptersCache(params domain.SearchParams) (domain.SearchResult, bool, error) {
	return domain.SearchResult{}, false, nil
}

func (c *NoopCacheRepository) PutChaptersCache(params domain.SearchParams, result domain.SearchResult) error {
	return nil
}
