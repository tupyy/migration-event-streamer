package elastic

import (
	"context"
	"fmt"
	"net/http"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/tupyy/migration-event-streamer/internal/clients"
	"github.com/tupyy/migration-event-streamer/internal/datastore/models"
	"go.uber.org/zap"
)

type ElasticDatasource struct {
	client      *elastic.Client
	bulkIndexer esutil.BulkIndexer
	index       string
	docIDPrefix string
}

func NewElasticDatastore(es *elastic.Client, config clients.ElasticSearchEnvConfig) (*ElasticDatasource, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         config.Index,           // The default index name
		Client:        es,                     // The Elasticsearch client
		NumWorkers:    config.NumWorkers,      // The number of worker goroutines
		FlushBytes:    int(config.FlushBytes), // The flush threshold in bytes
		FlushInterval: config.FlushInterval,   // The periodic flush interval
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the bulk indexer: %w", err)
	}

	elasticDt := &ElasticDatasource{
		client:      es,
		bulkIndexer: bi,
		index:       config.Index,
		docIDPrefix: config.DocIdPrefix,
	}

	// TODO find another place to do this.
	if err := elasticDt.createIndex(config.Index); err != nil {
		return nil, err
	}

	return elasticDt, nil
}

func (e *ElasticDatasource) Write(ctx context.Context, event models.Event) error {
	return e.bulkIndexer.Add(ctx, esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: event.ID,
		Body:       event.Body,
		OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
			zap.S().Debugf("item sucessfully added: %s", event.ID)
		},
		OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
			zap.S().Errorf("failed to add item to indexer: %w", err)
		},
	})
}

func (e *ElasticDatasource) createIndex(name string) error {
	// check if index exists
	res, err := e.client.Indices.Exists([]string{name})
	if err != nil {
		return fmt.Errorf("failed to check if index %s exists: %w", name, err)
	}
	if res.StatusCode == http.StatusOK {
		res.Body.Close()
		return nil
	}

	// create the index
	res, err = e.client.Indices.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", name, err)
	}
	if res.IsError() {
		return fmt.Errorf("failed to create index %s: %w", name, err)
	}
	res.Body.Close()
	return nil
}
