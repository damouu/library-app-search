package application

import (
	"errors"
	"testing"

	"library-app-search/internal/domain"
)

type fakeCacheRepository struct {
	cachedResult domain.SearchResult
	found        bool
	fetchErr     error
	putErr       error

	fetchCalled     bool
	putCalled       bool
	lastFetchParams domain.SearchParams
	lastPutParams   domain.SearchParams
}

func (f *fakeCacheRepository) FetchChaptersCache(params domain.SearchParams) (domain.SearchResult, bool, error) {
	f.fetchCalled = true
	f.lastFetchParams = params

	if f.fetchErr != nil {
		return domain.SearchResult{}, false, f.fetchErr
	}

	return f.cachedResult, f.found, nil
}

func (f *fakeCacheRepository) PutChaptersCache(params domain.SearchParams, result domain.SearchResult) error {
	f.putCalled = true
	f.lastPutParams = params
	return f.putErr
}

type fakeSearchRepository struct {
	result domain.SearchResult
	err    error

	searchCalled     bool
	lastSearchParams domain.SearchParams
}

func (f *fakeSearchRepository) SearchChapters(params domain.SearchParams) (domain.SearchResult, error) {
	f.searchCalled = true
	f.lastSearchParams = params
	if f.err != nil {
		return domain.SearchResult{}, f.err
	}

	return f.result, nil
}

func TestSearchService_PropagatesSearchParams(t *testing.T) {
	params := domain.SearchParams{
		Query: "one piece",
		Page:  2,
		Size:  3,
	}

	cacheRepository := &fakeCacheRepository{
		found: false,
	}

	searchRepository := &fakeSearchRepository{
		result: domain.SearchResult{
			Items: []domain.Chapter{},
			Page:  params.Page,
			Size:  params.Size,
			Total: 0,
		},
	}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters(params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cacheRepository.lastFetchParams != params {
		t.Fatalf(
			"unexpected cache fetch params: got %+v, want %+v",
			cacheRepository.lastFetchParams,
			params,
		)
	}

	if searchRepository.lastSearchParams != params {
		t.Fatalf(
			"unexpected search params: got %+v, want %+v",
			searchRepository.lastSearchParams,
			params,
		)
	}

	if cacheRepository.lastPutParams != params {
		t.Fatalf(
			"unexpected cache put params: got %+v, want %+v",
			cacheRepository.lastPutParams,
			params,
		)
	}
}

func TestSearchService_ReturnsCachedChapters(t *testing.T) {
	cachedChapters := []domain.Chapter{
		{
			ChapterUUID: "chapter-1",
			Title:       "One Piece",
		},
	}

	cacheRepository := &fakeCacheRepository{
		cachedResult: domain.SearchResult{
			Items: cachedChapters,
			Page:  1,
			Size:  10,
			Total: 1,
		},
		found: true,
	}

	searchRepository := &fakeSearchRepository{}

	service := NewSearchService(cacheRepository, searchRepository)

	result, err := service.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cacheRepository.fetchCalled {
		t.Fatal("expected cache repository to be called")
	}

	if searchRepository.searchCalled {
		t.Fatal("expected search repository not to be called when cache contains data")
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(result.Items))
	}

	if result.Items[0].Title != "One Piece" {
		t.Fatalf(
			"expected title %q, got %q",
			"One Piece",
			result.Items[0].Title,
		)
	}
}

func TestSearchService_SearchesElasticsearchOnCacheMiss(t *testing.T) {
	chapters := []domain.Chapter{
		{
			ChapterUUID: "chapter-1",
			Title:       "One Piece",
		},
	}

	cacheRepository := &fakeCacheRepository{
		found: false,
	}

	searchRepository := &fakeSearchRepository{
		result: domain.SearchResult{
			Items: chapters,
			Page:  1,
			Size:  10,
			Total: 1,
		},
	}

	service := NewSearchService(cacheRepository, searchRepository)

	result, err := service.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !searchRepository.searchCalled {
		t.Fatal("expected search repository to be called")
	}

	if !cacheRepository.putCalled {
		t.Fatal("expected results to be stored in cache")
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(result.Items))
	}
}

func TestSearchService_ReturnsErrorWhenCacheFails(t *testing.T) {
	expectedErr := errors.New("cache unavailable")

	cacheRepository := &fakeCacheRepository{
		fetchErr: expectedErr,
	}

	searchRepository := &fakeSearchRepository{}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if searchRepository.searchCalled {
		t.Fatal("expected search repository not to be called")
	}
}

func TestSearchService_ReturnsErrorWhenSearchFails(t *testing.T) {
	expectedErr := errors.New("elasticsearch unavailable")

	cacheRepository := &fakeCacheRepository{
		found: false,
	}

	searchRepository := &fakeSearchRepository{
		err: expectedErr,
	}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if !cacheRepository.putCalled {
		return
	}

	t.Fatal("expected cache repository not to be called after search failure")
}

func TestSearchService_ReturnsErrorWhenCacheWriteFails(t *testing.T) {
	expectedErr := errors.New("cache write failed")

	cacheRepository := &fakeCacheRepository{
		found:  false,
		putErr: expectedErr,
	}

	searchRepository := &fakeSearchRepository{
		result: domain.SearchResult{
			Items: []domain.Chapter{},
			Page:  1,
			Size:  10,
			Total: 0,
		},
	}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
