package elastic

import (
	"context"

	"github.com/opensearch-project/opensearch-go"
	"github.com/tupyy/migration-event-streamer/internal/datastore/models"
	"go.uber.org/zap"
)

type ElasticDatasource struct {
	client *opensearch.Client
}

func NewElasticDatastore(client *opensearch.Client) *ElasticDatasource {
	return &ElasticDatasource{client}
}

func (e *ElasticDatasource) Write(ctx context.Context, event models.Event) error {
	zap.S().Infof("******** %+v", event)
	return nil
}
