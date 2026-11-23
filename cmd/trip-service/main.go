package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/swokeqq/tripngo.git/internal/config"
	"github.com/swokeqq/tripngo.git/internal/container"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", slog.Any("error", err))
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	if err := run(cfg, logger); err != nil {
		logger.Error("application error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := container.New(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	srvErr := make(chan error, 1)
	go func() {
		if err := app.Run(); err != nil {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		if stopErr := shutdown(app, cfg.ShutdownTimeout); stopErr != nil {
			return fmt.Errorf("server startup failed: %w, stop failed too: %w", err, stopErr)
		}
		return fmt.Errorf("server startup failed: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	if err := shutdown(app, cfg.ShutdownTimeout); err != nil {
		return fmt.Errorf("forced shutdown: %w", err)
	}

	logger.Info("application stopped gracefully")
	return nil
}

func shutdown(app *container.Container, timeout time.Duration) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		return err
	}

	return nil
}
