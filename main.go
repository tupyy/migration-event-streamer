package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tupyy/migration-event-streamer/internal/clients"
	"github.com/tupyy/migration-event-streamer/internal/datastore"
	"github.com/tupyy/migration-event-streamer/internal/logger"
	"github.com/tupyy/migration-event-streamer/internal/services"
	"go.uber.org/zap"
)

func main() {
	logger := logger.SetupLogger()
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	pgConnection, err := clients.NewPgConnectionFromEnv()
	if err != nil {
		panic(err)
	}
	_, err = clients.NewElasticsearchClientFromEnv()
	if err != nil {
		panic(err)
	}

	dt := datastore.NewDatastore(pgConnection, nil)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	srv := services.NewElastic(dt, 2*time.Second)
	go srv.Run(ctx)

	<-ctx.Done()
}
