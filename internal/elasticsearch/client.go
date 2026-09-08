package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v9"
	"library-app-search/internal/config"
)

func CreateClient(cfg *config.ElasticsearchConfig) (*elasticsearch.Client, error) {

	es, err := elasticsearch.New(
		elasticsearch.WithAddresses(cfg.URL),
		elasticsearch.WithAPIKey(cfg.APIKey),
	)

	if err != nil {
		return nil, err
	}

	return es, nil
}
