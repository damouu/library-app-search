package handler

import (
	"encoding/json"
	"library-app-search/internal/domain"
	"net/http"
)

type SearchUseCase interface {
	SearchChapters(query string) ([]domain.Chapter, error)
}

type SearchHandler struct {
	searchUseCase SearchUseCase
}

func NewSearchHandler(searchUseCase SearchUseCase) *SearchHandler {
	return &SearchHandler{searchUseCase: searchUseCase}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "400 Bad Request query parameter is required", http.StatusBadRequest)
		return
	}
	chapters, err := s.searchUseCase.SearchChapters(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(chapters)
	if err != nil {
		return
	}
}
