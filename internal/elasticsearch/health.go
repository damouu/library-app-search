package elasticsearch

import "github.com/elastic/go-elasticsearch/v9"

type HealthChecker struct {
	client *elasticsearch.Client
}

func NewHealthChecker(client *elasticsearch.Client) *HealthChecker {
	return &HealthChecker{
		client: client,
	}
}

func (h *HealthChecker) Check() error {
	_, err := h.client.Info()
	return err
}
