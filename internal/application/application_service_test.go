package application

import (
	"errors"
	"testing"

	"library-app-search/internal/domain"
)

type fakeCacheRepository struct {
	cachedChapters []domain.Chapter
	found          bool
	fetchErr       error
	putErr         error

	fetchCalled bool
	putCalled   bool
}

func (f *fakeCacheRepository) FetchChaptersCache(query string) ([]domain.Chapter, bool, error) {
	f.fetchCalled = true

	if f.fetchErr != nil {
		return nil, false, f.fetchErr
	}

	return f.cachedChapters, f.found, nil
}

func (f *fakeCacheRepository) PutChaptersCache(query string, chapters []domain.Chapter) error {
	f.putCalled = true

	return f.putErr
}

type fakeSearchRepository struct {
	chapters []domain.Chapter
	err      error

	searchCalled bool
}

func (f *fakeSearchRepository) SearchChapters(query string) ([]domain.Chapter, error) {
	f.searchCalled = true

	if f.err != nil {
		return nil, f.err
	}

	return f.chapters, nil
}

func TestSearchService_ReturnsCachedChapters(t *testing.T) {
	cachedChapters := []domain.Chapter{
		{
			ChapterUUID: "chapter-1",
			Title:       "One Piece",
		},
	}

	cacheRepository := &fakeCacheRepository{
		cachedChapters: cachedChapters,
		found:          true,
	}

	searchRepository := &fakeSearchRepository{}

	service := NewSearchService(cacheRepository, searchRepository)

	result, err := service.SearchChapters("one piece")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cacheRepository.fetchCalled {
		t.Fatal("expected cache repository to be called")
	}

	if searchRepository.searchCalled {
		t.Fatal("expected search repository not to be called when cache contains data")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(result))
	}

	if result[0].Title != "One Piece" {
		t.Fatalf("expected title %q, got %q", "One Piece", result[0].Title)
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
		chapters: chapters,
	}

	service := NewSearchService(cacheRepository, searchRepository)

	result, err := service.SearchChapters("one piece")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !searchRepository.searchCalled {
		t.Fatal("expected search repository to be called")
	}

	if !cacheRepository.putCalled {
		t.Fatal("expected results to be stored in cache")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(result))
	}
}

func TestSearchService_ReturnsErrorWhenCacheFails(t *testing.T) {
	expectedErr := errors.New("cache unavailable")

	cacheRepository := &fakeCacheRepository{
		fetchErr: expectedErr,
	}

	searchRepository := &fakeSearchRepository{}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters("one piece")

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

	_, err := service.SearchChapters("one piece")

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
		chapters: []domain.Chapter{},
	}

	service := NewSearchService(cacheRepository, searchRepository)

	_, err := service.SearchChapters("one piece")

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
