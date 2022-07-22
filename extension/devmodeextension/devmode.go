package devmode

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var _ component.Extension = (*devMode)(nil)

type devMode struct {
	logger *zap.Logger
}

func newDevMode(c *Config, logger *zap.Logger) (component.Extension, error) {
	return &devMode{
		logger: logger,
	}, nil
}

func (d devMode) Start(ctx context.Context, host component.Host) error {
	d.logger.Info("starting devmode!")
	err := startServer(context.Background())

	return err
}

func (d devMode) Shutdown(ctx context.Context) error {
	d.logger.Info("shutting down devmode!")

	return nil
}
