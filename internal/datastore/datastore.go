package datastore

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/tupyy/migration-event-streamer/internal/clients"
	"github.com/tupyy/migration-event-streamer/internal/datastore/elastic"
	"github.com/tupyy/migration-event-streamer/internal/datastore/postgres"
	"github.com/tupyy/migration-event-streamer/pkg/datastore"
)

type pgElasticDatastore struct {
	reader *postgres.PgDatastore
	writer *elastic.ElasticDatasource
}

func NewDatastore(pgConn *clients.PgConnection, elasticConn *es.Client, elasticConf clients.ElasticSearchEnvConfig) (*pgElasticDatastore, error) {
	writer, err := elastic.NewElasticDatastore(elasticConn, elasticConf)
	if err != nil {
		return nil, err
	}
	return &pgElasticDatastore{
		reader: postgres.NewPgDatastore(pgConn),
		writer: writer,
	}, nil
}

func (d *pgElasticDatastore) ReadWriteTx(ctx context.Context, fn func(context.Context, datastore.Reader, datastore.Writer) error) error {
	return fn(ctx, d.reader, d.writer)
}
