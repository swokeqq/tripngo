package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/swokeqq/tripngo.git/internal/app"
	"github.com/swokeqq/tripngo.git/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", slog.Any("error", err))
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, logger); err != nil {
		logger.Error("application error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	srvErr := make(chan error, 1)
	go func() {
		if err := application.Run(); err != nil {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		return fmt.Errorf("server startup failed: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := application.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("forced shutdown: %w", err)
	}

	logger.Info("application stopped gracefully")
	return nil
}
