package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"library-app-search/internal/domain"
)

// SearchHit represents an OpenSearch search hit.
type SearchHit struct {
	Source domain.Chapter `json:"_source"`
}

// SearchHits contains the total number of matches and the returned documents.
type SearchHits struct {
	Total TotalHits   `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

// TotalHits contains the total number of matching documents.
type TotalHits struct {
	Value int `json:"value"`
}

// SearchResponse represents the relevant part of the OpenSearch search response.
type SearchResponse struct {
	Hits SearchHits `json:"hits"`
}

// SearchRepository provides chapter search operations using OpenSearch.
type SearchRepository struct {
	client *opensearchapi.Client
}

// NewSearchRepository creates a search repository using the provided OpenSearch client.
func NewSearchRepository(client *opensearchapi.Client) *SearchRepository {
	return &SearchRepository{
		client: client,
	}
}

// SearchChapters searches chapters by title and second title with pagination.
func (s *SearchRepository) SearchChapters(params domain.SearchParams) (domain.SearchResult, error) {
	searchQuery := SearchQuery{
		From: (params.Page - 1) * params.Size,
		Size: params.Size,
		Query: QueryClause{
			MultiMatch: MultiMatch{
				Query:  params.Query,
				Fields: []string{"title", "second_title"},
			},
		},
	}

	requestBody, err := json.Marshal(searchQuery)
	if err != nil {
		return domain.SearchResult{}, err
	}

	response, err := s.client.Search(
		context.Background(),
		&opensearchapi.SearchReq{
			Indices: []string{"chapters"},
			Body:    bytes.NewReader(requestBody),
		},
	)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("search OpenSearch: %w", err)
	}

	chapters := make([]domain.Chapter, 0, len(response.Hits.Hits))

	for _, hit := range response.Hits.Hits {
		var chapter domain.Chapter

		if err := json.Unmarshal(hit.Source, &chapter); err != nil {
			return domain.SearchResult{}, fmt.Errorf("decode chapter: %w", err)
		}

		chapters = append(chapters, chapter)
	}

	return domain.SearchResult{
		Items: chapters,
		Page:  params.Page,
		Size:  params.Size,
		Total: response.Hits.Total.Value,
	}, nil
}
