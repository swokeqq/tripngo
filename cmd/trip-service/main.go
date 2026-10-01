package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/swokeqq/tripngo.git/internal/app"
	"github.com/swokeqq/tripngo.git/internal/config"
	"github.com/swokeqq/tripngo.git/internal/handler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	h := handler.NewHandler(nil)
	application := app.New(cfg, h)

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
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := application.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("forced shutdown: %v", err)
	}

	log.Println("application stopped gracefully")
	return nil
}
