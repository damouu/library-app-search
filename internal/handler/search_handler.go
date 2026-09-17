package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"library-app-search/internal/domain"
)

// SearchUseCase defines the application operation required by the search handler.
type SearchUseCase interface {
	SearchChapters(params domain.SearchParams) (domain.SearchResult, error)
}

// SearchHandler handles HTTP requests for chapter searches.
type SearchHandler struct {
	searchUseCase SearchUseCase
}

// NewSearchHandler creates a new search handler with the required use case.
func NewSearchHandler(searchUseCase SearchUseCase) *SearchHandler {
	return &SearchHandler{
		searchUseCase: searchUseCase,
	}
}

// Search handles chapter search requests and validates pagination parameters.
func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Reject empty or whitespace-only search queries.
	if query == "" {
		http.Error(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	// Use sensible defaults when pagination parameters are omitted.
	page := 1
	size := 10

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err != nil || parsedPage < 1 || parsedPage > 1000 {
			http.Error(
				w,
				"page must be an integer between 1 and 1000",
				http.StatusBadRequest,
			)
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

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
