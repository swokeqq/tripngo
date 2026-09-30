package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr        string        `env:"HTTP_ADDR,required"`
	DatabaseURL     string        `env:"DATABASE_URL,required"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`

	// Database
	DatabaseMaxConns  int32         `env:"DATABASE_MAX_CONNS" envDefault:"10"`
	DatabaseMinConns  int32         `env:"DATABASE_MIN_CONNS" envDefault:"2"`
	DBMaxConnLifetime time.Duration `env:"DATABASE_MAX_CONN_LIFETIME" envDefault:"30m"`
	DBConnectTimeout  time.Duration `env:"DATABASE_CONNECT_TIMEOUT" envDefault:"5s"`
	DBQueryTimeout    time.Duration `env:"DATABASE_QUERY_TIMEOUT" envDefault:"3s"`

	// HTTP Timeouts
	HTTPReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"2s"`
	HTTPReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	HTTPWriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	HTTPIdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"30s"`
}

func Load() (*Config, error) {

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.DatabaseMaxConns <= 0 {
		return errors.New("DATABASE_MAX_CONNS must be greater than zero")
	}

	if cfg.DatabaseMinConns < 0 {
		return errors.New("DATABASE_MIN_CONNS cannot be negative")
	}

	if cfg.DatabaseMaxConns < cfg.DatabaseMinConns {
		return fmt.Errorf("DATABASE_MIN_CONNS (%d) cannot be greater than DATABASE_MAX_CONNS (%d)",
			cfg.DatabaseMinConns, cfg.DatabaseMaxConns)
	}

	durations := map[string]time.Duration{
		"SHUTDOWN_TIMEOUT":           cfg.ShutdownTimeout,
		"DATABASE_MAX_CONN_LIFETIME": cfg.DBMaxConnLifetime,
		"DATABASE_CONNECT_TIMEOUT":   cfg.DBConnectTimeout,
		"DATABASE_QUERY_TIMEOUT":     cfg.DBQueryTimeout,
		"HTTP_READ_HEADER_TIMEOUT":   cfg.HTTPReadHeaderTimeout,
		"HTTP_READ_TIMEOUT":          cfg.HTTPReadTimeout,
		"HTTP_WRITE_TIMEOUT":         cfg.HTTPWriteTimeout,
		"HTTP_IDLE_TIMEOUT":          cfg.HTTPIdleTimeout,
	}

	for name, dur := range durations {
		if dur <= 0 {
			return fmt.Errorf("%s must be positive, got: %v", name, dur)
		}
	}

	return nil
}
