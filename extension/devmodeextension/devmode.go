package devmodeextension

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var Storage Storer

var _ component.Extension = (*devMode)(nil)

type devMode struct {
	logger  *zap.Logger
	storage *dbStorageClient
}

func newDevMode(ctx context.Context, c *Config, logger *zap.Logger) (component.Extension, error) {
	client, err := newClient(ctx, "sqlite3", "spans", logger)
	if err != nil {
		return nil, err
	}

	Storage = client

	return &devMode{
		logger:  logger,
		storage: client,
	}, nil
}

func (d *devMode) Start(ctx context.Context, host component.Host) error {
	d.logger.Info("starting devmode!")
	err := startServer(context.Background(), d.logger, host)

	return err
}

func (d *devMode) Shutdown(ctx context.Context) error {
	d.logger.Info("shutting down devmode!")

	return nil
}
