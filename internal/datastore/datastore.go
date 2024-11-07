package datastore

import (
	"context"

	"github.com/opensearch-project/opensearch-go"
	"github.com/tupyy/migration-event-streamer/internal/clients"
	"github.com/tupyy/migration-event-streamer/internal/datastore/elastic"
	"github.com/tupyy/migration-event-streamer/internal/datastore/models"
	"github.com/tupyy/migration-event-streamer/internal/datastore/postgres"
)

type Reader interface {
	Read(context.Context) ([]models.Source, error)
}

type Writer interface {
	Write(context.Context, models.Event) error
}

type Datastore interface {
	ReadWriteTx(context.Context, func(context.Context, Reader, Writer) error) error
}

type pgElasticDatastore struct {
	reader *postgres.PgDatastore
	writer *elastic.ElasticDatasource
}

func NewDatastore(pgConn *clients.PgConnection, elasticConn *opensearch.Client) *pgElasticDatastore {
	return &pgElasticDatastore{
		reader: postgres.NewPgDatastore(pgConn),
		writer: elastic.NewElasticDatastore(elasticConn),
	}
}

func (d *pgElasticDatastore) ReadWriteTx(ctx context.Context, fn func(context.Context, Reader, Writer) error) error {
	return fn(ctx, d.reader, d.writer)
}
