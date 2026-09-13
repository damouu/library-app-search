package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"library-app-search/internal/domain"
)

type SearchHit struct {
	Source domain.Chapter `json:"_source"`
}

type SearchHits struct {
	Hits []SearchHit `json:"hits"`
}

type SearchResponse struct {
	Hits SearchHits `json:"hits"`
}

type SearchRepository struct {
	client *elasticsearch.Client
}

func NewSearchRepository(client *elasticsearch.Client) *SearchRepository {
	return &SearchRepository{
		client: client,
	}
}

func (s *SearchRepository) SearchChapters(query string) ([]domain.Chapter, error) {
	searchQuery := SearchQuery{
		Query: QueryClause{MultiMatch: MultiMatch{
			Query:  query,
			Fields: []string{"title", "second_title", "summary"},
		}},
	}
	jsonData, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}
	requestBody := bytes.NewReader(jsonData)

	response, err := s.client.Search(
		s.client.Search.WithIndex("chapters"),
		s.client.Search.WithBody(requestBody),
	)

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", response.Status())
	}

	defer response.Body.Close()

	var searchResponse SearchResponse

	err = json.NewDecoder(response.Body).Decode(&searchResponse)
	if err != nil {
		return nil, err
	}

	chapters := make([]domain.Chapter, 0)

	for _, hit := range searchResponse.Hits.Hits {
		chapters = append(chapters, hit.Source)
	}
	return chapters, nil

}
