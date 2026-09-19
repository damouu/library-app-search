package opensearch

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type HealthChecker struct {
	client *opensearchapi.Client
}

func NewHealthChecker(client *opensearchapi.Client) *HealthChecker {
	return &HealthChecker{
		client: client,
	}
}

func (h *HealthChecker) Check() error {
	_, err := h.client.Info(context.Background(), nil)
	return err
}
