package opensearch

import (
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"library-app-search/internal/config"
)

func CreateClient(cfg *config.OpenSearchConfig) (*opensearchapi.Client, error) {
	return opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearch.Config{
				Addresses: []string{cfg.URL},
				Username:  cfg.Username,
				Password:  cfg.Password,
			},
		},
	)
}
