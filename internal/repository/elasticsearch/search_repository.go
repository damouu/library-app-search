package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
	"library-app-search/internal/domain"
)

// SearchHit represents an Elasticsearch search hit.
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

// SearchResponse represents the relevant part of the Elasticsearch search response.
type SearchResponse struct {
	Hits SearchHits `json:"hits"`
}

// SearchRepository provides chapter search operations using Elasticsearch.
type SearchRepository struct {
	client *elasticsearch.Client
}

// NewSearchRepository creates a search repository using the provided Elasticsearch client.
func NewSearchRepository(client *elasticsearch.Client) *SearchRepository {
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
		s.client.Search.WithIndex("chapters"),
		s.client.Search.WithBody(bytes.NewReader(requestBody)),
	)
	if err != nil {
		return domain.SearchResult{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		return domain.SearchResult{}, fmt.Errorf(
			"elasticsearch error: %s",
			response.Status(),
		)
	}

	var searchResponse SearchResponse

	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return domain.SearchResult{}, err
	}

	chapters := make([]domain.Chapter, 0, len(searchResponse.Hits.Hits))

	for _, hit := range searchResponse.Hits.Hits {
		chapters = append(chapters, hit.Source)
	}

	return domain.SearchResult{
		Items: chapters,
		Page:  params.Page,
		Size:  params.Size,
		Total: searchResponse.Hits.Total.Value,
	}, nil
}
