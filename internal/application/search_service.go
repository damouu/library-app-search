package application

import "library-app-search/internal/domain"

type CacheRepository interface {
	FetchChaptersCache(params domain.SearchParams) (domain.SearchResult, bool, error)
	PutChaptersCache(params domain.SearchParams, result domain.SearchResult) error
}

type SearchRepository interface {
	SearchChapters(params domain.SearchParams) (domain.SearchResult, error)
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

func (s *SearchService) SearchChapters(params domain.SearchParams) (domain.SearchResult, error) {
	result, found, err := s.cacheRepository.FetchChaptersCache(params)
	if err != nil {
		return domain.SearchResult{}, err
	}

	if found {
		return result, nil
	}

	result, err = s.searchRepository.SearchChapters(params)
	if err != nil {
		return domain.SearchResult{}, err
	}

	err = s.cacheRepository.PutChaptersCache(params, result)
	if err != nil {
		return domain.SearchResult{}, err
	}

	return result, nil
}
