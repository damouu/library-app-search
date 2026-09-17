package elasticsearch

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"library-app-search/internal/domain"
)

func TestSearchRepository_SearchChapters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/chapters/_search" {
			t.Fatalf("expected /chapters/_search, got %s", r.URL.Path)
		}

		var actualQuery SearchQuery

		if err := json.NewDecoder(r.Body).Decode(&actualQuery); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		expectedQuery := SearchQuery{
			From: 3,
			Size: 3,
			Query: QueryClause{
				MultiMatch: MultiMatch{
					Query:  "one piece",
					Fields: []string{"title", "second_title"},
				},
			},
		}

		if !reflect.DeepEqual(actualQuery, expectedQuery) {
			t.Fatalf(
				"unexpected search query:\ngot:  %+v\nwant: %+v",
				actualQuery,
				expectedQuery,
			)
		}

		response := map[string]interface{}{
			"hits": map[string]interface{}{
				"total": map[string]interface{}{
					"value": 2,
				},
				"hits": []map[string]interface{}{
					{
						"_source": map[string]interface{}{
							"chapter_uuid":      "chapter-1",
							"series_uuid":       "series-1",
							"title":             "One Piece",
							"second_title":      "ONE PIECE",
							"summary":           "A pirate adventure",
							"chapter_number":    15,
							"total_pages":       177,
							"publication_date":  "2026-07-10",
							"cover_artwork_url": "https://example.com/cover.jpg",
						},
					},
					{
						"_source": map[string]interface{}{
							"chapter_uuid":      "chapter-2",
							"series_uuid":       "series-1",
							"title":             "One Piece 2",
							"second_title":      "ONE PIECE 2",
							"summary":           "Another pirate adventure",
							"chapter_number":    16,
							"total_pages":       177,
							"publication_date":  "2026-07-11",
							"cover_artwork_url": "https://example.com/cover-2.jpg",
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{server.URL},
	})
	if err != nil {
		t.Fatalf("failed to create Elasticsearch client: %v", err)
	}

	repository := NewSearchRepository(client)

	params := domain.SearchParams{
		Query: "one piece",
		Page:  2,
		Size:  3,
	}

	result, err := repository.SearchChapters(params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != params.Page {
		t.Fatalf("expected page %d, got %d", params.Page, result.Page)
	}

	if result.Size != params.Size {
		t.Fatalf("expected size %d, got %d", params.Size, result.Size)
	}

	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 chapters, got %d", len(result.Items))
	}

	if result.Items[0].ChapterUUID != "chapter-1" {
		t.Fatalf(
			"expected first chapter UUID %q, got %q",
			"chapter-1",
			result.Items[0].ChapterUUID,
		)
	}

	if result.Items[0].Title != "One Piece" {
		t.Fatalf(
			"expected first title %q, got %q",
			"One Piece",
			result.Items[0].Title,
		)
	}

	if result.Items[1].ChapterUUID != "chapter-2" {
		t.Fatalf(
			"expected second chapter UUID %q, got %q",
			"chapter-2",
			result.Items[1].ChapterUUID,
		)
	}
}

func TestSearchRepository_ReturnsErrorWhenElasticsearchFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = w.Write([]byte(`{
			"error": {
				"reason": "search service unavailable"
			}
		}`))
	}))
	defer server.Close()

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{server.URL},
	})
	if err != nil {
		t.Fatalf("failed to create Elasticsearch client: %v", err)
	}

	repository := NewSearchRepository(client)

	_, err = repository.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if err == nil {
		t.Fatal("expected Elasticsearch error, got nil")
	}

	if !strings.Contains(err.Error(), "500 Internal Server Error") {
		t.Fatalf(
			"expected error to contain %q, got %q",
			"500 Internal Server Error",
			err.Error(),
		)
	}
}

func TestSearchRepository_ReturnsErrorWhenResponseIsInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{server.URL},
	})
	if err != nil {
		t.Fatalf("failed to create Elasticsearch client: %v", err)
	}

	repository := NewSearchRepository(client)

	_, err = repository.SearchChapters(domain.SearchParams{
		Query: "one piece",
		Page:  1,
		Size:  10,
	})

	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}
}
