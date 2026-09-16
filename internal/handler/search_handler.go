package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"library-app-search/internal/domain"
)

type SearchUseCase interface {
	SearchChapters(params domain.SearchParams) (domain.SearchResult, error)
}

type SearchHandler struct {
	searchUseCase SearchUseCase
}

func NewSearchHandler(searchUseCase SearchUseCase) *SearchHandler {
	return &SearchHandler{
		searchUseCase: searchUseCase,
	}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	page := 1
	size := 10

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err != nil || parsedPage < 1 {
			http.Error(w, "page must be a positive integer", http.StatusBadRequest)
			return
		}

		page = parsedPage
	}

	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		parsedSize, err := strconv.Atoi(sizeParam)
		if err != nil || parsedSize < 1 || parsedSize > 100 {
			http.Error(
				w,
				"size must be an integer between 1 and 100",
				http.StatusBadRequest,
			)
			return
		}

		size = parsedSize
	}

	params := domain.SearchParams{
		Query: query,
		Page:  page,
		Size:  size,
	}

	result, err := s.searchUseCase.SearchChapters(params)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
