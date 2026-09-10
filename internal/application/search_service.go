package application

import "library-app-search/internal/domain"

type CacheRepository interface {
	FetchChaptersCache(query string) ([]domain.Chapter, bool, error)
	PutChaptersCache(query string, chapters []domain.Chapter) error
}

type SearchRepository interface {
	SearchChapters(query string) ([]domain.Chapter, error)
}

type SearchService struct {
	cacheRepository  CacheRepository
	searchRepository SearchRepository
}

func NewSearchService(cacheRepository CacheRepository, searchRepository SearchRepository) *SearchService {
	return &SearchService{
		cacheRepository:  cacheRepository,
		searchRepository: searchRepository,
	}
}

func (s *SearchService) SearchChapters(query string) ([]domain.Chapter, error) {
	chapters, found, err := s.cacheRepository.FetchChaptersCache(query)
	if err != nil {
		return nil, err
	}

	if found {
		return chapters, nil
	}

	chapters, err = s.searchRepository.SearchChapters(query)
	if err != nil {
		return nil, err
	}

	err = s.cacheRepository.PutChaptersCache(query, chapters)
	if err != nil {
		return nil, err
	}

	return chapters, nil
}
