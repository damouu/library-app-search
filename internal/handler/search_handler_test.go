package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-app-search/internal/domain"
)

type fakeSearchUseCase struct {
	chapters []domain.Chapter
	err      error

	called bool
}

func (f *fakeSearchUseCase) SearchChapters(query string) ([]domain.Chapter, error) {
	f.called = true

	if f.err != nil {
		return nil, f.err
	}

	return f.chapters, nil
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
		chapters: []domain.Chapter{
			{
				ChapterUUID: "chapter-1",
				Title:       "One Piece",
			},
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

	expectedBody := `[{"chapter_uuid":"chapter-1","series_uuid":"","title":"One Piece","second_title":"","summary":"","chapter_number":0,"total_pages":0,"publication_date":"","cover_artwork_url":""}]` + "\n"

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
		chapters: []domain.Chapter{},
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

	expectedBody := "[]\n"

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
