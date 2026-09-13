package redis

import "library-app-search/internal/domain"

type NoopCacheRepository struct{}

func NewNoopCacheRepository() *NoopCacheRepository {
	return &NoopCacheRepository{}
}

func (c *NoopCacheRepository) FetchChaptersCache(query string) ([]domain.Chapter, bool, error) {
	return nil, false, nil
}

func (c *NoopCacheRepository) PutChaptersCache(query string, chapters []domain.Chapter) error {
	return nil
}
