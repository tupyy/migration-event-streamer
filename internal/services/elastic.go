package services

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tupyy/migration-event-streamer/internal/datastore/models"
	"github.com/tupyy/migration-event-streamer/internal/transform"
	"github.com/tupyy/migration-event-streamer/pkg/datastore"
	"go.uber.org/zap"
)

type Inventory struct {
	dt          datastore.Datastore
	readTimeout time.Duration
}

func NewInventory(dt datastore.Datastore, readTimeout time.Duration) *Inventory {
	return &Inventory{dt: dt, readTimeout: readTimeout}
}

func (e *Inventory) Run(ctx context.Context) {
	zap.S().Infof("start inserting events every %s", e.readTimeout)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		<-time.After(e.readTimeout)

		if err := e.dt.ReadWriteTx(ctx, func(ctx context.Context, reader datastore.Reader, writer datastore.Writer) error {
			// read
			sources, err := reader.Read(ctx)
			if err != nil {
				return err
			}

			if len(sources) == 0 {
				return nil
			}

			for _, source := range sources {
				// transform the source inventory to elastic inventory
				inventory := transform.InventorySourceToElastic(source)

				// marshal and create the event
				data, err := json.Marshal(inventory)
				if err != nil {
					zap.S().Warnw("failed to marshal inventory", "error", err, "inventory", inventory)
					continue
				}

				event := models.Event{
					Index: "assited-migrations",
					ID:    uuid.New().String(),
					Body:  bytes.NewReader(data),
				}

				// write
				if err := writer.Write(ctx, event); err != nil {
					zap.S().Warnw("failed to write inventory", "error", err, "inventory", inventory)
					continue
				}
			}

			return nil
		}); err != nil {
			zap.S().Errorf("failed to write inventory to elastic: %s", err)
		}
	}
}
