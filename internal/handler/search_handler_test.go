package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-app-search/internal/domain"
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
					ChapterUUID: "chapter-1",
					Title:       "One Piece",
				},
			},
			Page:  1,
			Size:  10,
			Total: 1,
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

	expectedBody := `{"items":[{"chapter_uuid":"chapter-1","series_uuid":"","title":"One Piece","second_title":"","summary":"","chapter_number":0,"total_pages":0,"publication_date":"","cover_artwork_url":""}],"page":1,"size":10,"total":1}` + "\n"

	if recorder.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %s, got %s",
			expectedBody,
			recorder.Body.String(),
		)
	}
}

func TestSearchHandler_ReturnsEmptyArrayWhenNoChapterIsFound(t *testing.T) {
	useCase := &fakeSearchUseCase{
		result: domain.SearchResult{
			Items: []domain.Chapter{},
			Page:  1,
			Size:  10,
			Total: 0,
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

	expectedBody := `{"items":[],"page":1,"size":10,"total":0}` + "\n"

	if recorder.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %s, got %s",
			expectedBody,
			recorder.Body.String(),
		)
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
			Items: []domain.Chapter{},
			Page:  2,
			Size:  3,
			Total: 13,
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
