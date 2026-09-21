package application

import (
	"errors"
	"reflect"
	"testing"
	uuid2 "uuid"

	"github.com/google/uuid"

	"library-app-search/internal/domain"
)

var (
	testChapterUUID = uuid.MustParse("ea0fbdc6-8b84-4c42-9ce8-8e07d9c29818")
	testSeriesUUID  = uuid.MustParse("df1b59e5-d788-4e5b-975a-161ec33e1d0e")
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
	lastPutResult   domain.SearchResult
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
	f.lastPutResult = result

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
			ChapterUUID: uuid2.UUID(testChapterUUID),
			SeriesUUID:  uuid2.UUID(testSeriesUUID),
			Title:       "One Piece",
		},
	}

	cachedResult := domain.SearchResult{
		Items:      cachedChapters,
		Page:       1,
		Size:       10,
		Total:      1,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	cacheRepository := &fakeCacheRepository{
		cachedResult: cachedResult,
		found:        true,
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

	if cacheRepository.putCalled {
		t.Fatal("expected cache repository not to be called when cache contains data")
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

func TestSearchService_SearchesOpenSearchOnCacheMiss(t *testing.T) {
	chapters := []domain.Chapter{
		{
			ChapterUUID: uuid2.UUID(testChapterUUID),
			SeriesUUID:  uuid2.UUID(testSeriesUUID),
			Title:       "One Piece",
		},
	}

	searchResult := domain.SearchResult{
		Items:      chapters,
		Page:       1,
		Size:       10,
		Total:      1,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	cacheRepository := &fakeCacheRepository{
		found: false,
	}

	searchRepository := &fakeSearchRepository{
		result: searchResult,
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

	if !reflect.DeepEqual(cacheRepository.lastPutResult, searchResult) {
		t.Fatalf(
			"unexpected cache result: got %+v, want %+v",
			cacheRepository.lastPutResult,
			searchResult,
		)
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
	expectedErr := errors.New("opensearch unavailable")

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

	if cacheRepository.putCalled {
		t.Fatal("expected cache repository not to be called after search failure")
	}
}

func TestSearchService_ReturnsErrorWhenCacheWriteFails(t *testing.T) {
	expectedErr := errors.New("cache write failed")

	cacheRepository := &fakeCacheRepository{
		found:  false,
		putErr: expectedErr,
	}

	searchRepository := &fakeSearchRepository{
		result: domain.SearchResult{
			Items:      []domain.Chapter{},
			Page:       1,
			Size:       10,
			Total:      0,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
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
