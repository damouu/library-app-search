package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	uuid2 "uuid"

	"github.com/google/uuid"

	"library-app-search/internal/domain"
)

var (
	testChapterUUID = uuid.MustParse("ea0fbdc6-8b84-4c42-9ce8-8e07d9c29818")
	testSeriesUUID  = uuid.MustParse("df1b59e5-d788-4e5b-975a-161ec33e1d0e")
	testDate        = time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC)
)

type fakeSearchUseCase struct {
	result domain.SearchResult
	err    error

	called bool
	params domain.SearchParams
}

func (f *fakeSearchUseCase) SearchChapters(params domain.SearchParams) (domain.SearchResult, error) {
	f.called = true
	f.params = params

	if f.err != nil {
		return domain.SearchResult{}, f.err
	}

	return f.result, nil
}

func TestSearchHandler_ReturnsBadRequestWhenQueryIsMissing(t *testing.T) {
	useCase := &fakeSearchUseCase{}
	searchHandler := NewSearchHandler(useCase)

	request := httptest.NewRequest(http.MethodGet, "/search", nil)
	recorder := httptest.NewRecorder()

	searchHandler.Search(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if useCase.called {
		t.Fatal("expected use case not to be called")
	}
}

func TestSearchHandler_ReturnsChapters(t *testing.T) {
	useCase := &fakeSearchUseCase{
		result: domain.SearchResult{
			Items: []domain.Chapter{
				{
					ChapterUUID:     uuid2.UUID(testChapterUUID),
					SeriesUUID:      uuid2.UUID(testSeriesUUID),
					Title:           "One Piece",
					SecondTitle:     "One Piece",
					Summary:         "Test summary",
					ChapterNumber:   1,
					TotalPages:      177,
					PublicationDate: testDate,
					CoverArtworkURL: "https://example.com/cover.jpg",
				},
			},
			Page:       1,
			Size:       10,
			Total:      1,
			TotalPages: 1,
			HasNext:    false,
			HasPrev:    false,
		},
	}

	searchHandler := NewSearchHandler(useCase)

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?q=one+piece",
		nil,
	)

	recorder := httptest.NewRecorder()

	searchHandler.Search(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	if !useCase.called {
		t.Fatal("expected use case to be called")
	}

	var response SearchResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(response.Items) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(response.Items))
	}

	chapter := response.Items[0]

	if chapter.ChapterUUID != testChapterUUID.String() {
		t.Fatalf(
			"expected chapter UUID %q, got %q",
			testChapterUUID.String(),
			chapter.ChapterUUID,
		)
	}

	if chapter.SeriesUUID != testSeriesUUID.String() {
		t.Fatalf(
			"expected series UUID %q, got %q",
			testSeriesUUID.String(),
			chapter.SeriesUUID,
		)
	}

	if chapter.Title != "One Piece" {
		t.Fatalf(
			"expected title %q, got %q",
			"One Piece",
			chapter.Title,
		)
	}

	if chapter.PublicationDate != "2026-07-10" {
		t.Fatalf(
			"expected publication date %q, got %q",
			"2026-07-10",
			chapter.PublicationDate,
		)
	}

	if response.Page != 1 {
		t.Fatalf("expected page %d, got %d", 1, response.Page)
	}

	if response.Size != 10 {
		t.Fatalf("expected size %d, got %d", 10, response.Size)
	}

	if response.Total != 1 {
		t.Fatalf("expected total %d, got %d", 1, response.Total)
	}

	if response.TotalPages != 1 {
		t.Fatalf("expected total pages %d, got %d", 1, response.TotalPages)
	}

	if response.HasNext {
		t.Fatal("expected has_next to be false")
	}

	if response.HasPrev {
		t.Fatal("expected has_prev to be false")
	}
}

func TestSearchHandler_ReturnsEmptyArrayWhenNoChapterIsFound(t *testing.T) {
	useCase := &fakeSearchUseCase{
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

	searchHandler := NewSearchHandler(useCase)

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?q=harry",
		nil,
	)

	recorder := httptest.NewRecorder()

	searchHandler.Search(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response SearchResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.Items == nil {
		t.Fatal("expected items to be an empty array, got nil")
	}

	if len(response.Items) != 0 {
		t.Fatalf("expected 0 chapters, got %d", len(response.Items))
	}

	if response.Page != 1 {
		t.Fatalf("expected page %d, got %d", 1, response.Page)
	}

	if response.Size != 10 {
		t.Fatalf("expected size %d, got %d", 10, response.Size)
	}

	if response.Total != 0 {
		t.Fatalf("expected total %d, got %d", 0, response.Total)
	}

	if response.TotalPages != 0 {
		t.Fatalf("expected total pages %d, got %d", 0, response.TotalPages)
	}
}

func TestSearchHandler_ReturnsInternalServerErrorWhenUseCaseFails(t *testing.T) {
	useCase := &fakeSearchUseCase{
		err: errors.New("search failed"),
	}

	searchHandler := NewSearchHandler(useCase)

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?q=one+piece",
		nil,
	)

	recorder := httptest.NewRecorder()

	searchHandler.Search(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}

func TestSearchHandler_PassesPaginationParametersToUseCase(t *testing.T) {
	useCase := &fakeSearchUseCase{
		result: domain.SearchResult{
			Items:      []domain.Chapter{},
			Page:       2,
			Size:       3,
			Total:      13,
			TotalPages: 5,
			HasNext:    true,
			HasPrev:    true,
		},
	}

	searchHandler := NewSearchHandler(useCase)

	request := httptest.NewRequest(
		http.MethodGet,
		"/search?q=one+piece&page=2&size=3",
		nil,
	)

	recorder := httptest.NewRecorder()

	searchHandler.Search(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	if useCase.params.Query != "one piece" {
		t.Fatalf(
			"expected query %q, got %q",
			"one piece",
			useCase.params.Query,
		)
	}

	if useCase.params.Page != 2 {
		t.Fatalf(
			"expected page %d, got %d",
			2,
			useCase.params.Page,
		)
	}

	if useCase.params.Size != 3 {
		t.Fatalf(
			"expected size %d, got %d",
			3,
			useCase.params.Size,
		)
	}
}
