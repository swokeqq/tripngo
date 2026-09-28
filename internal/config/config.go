package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	LogLevel    string

	ShutdownTimeout time.Duration

	DatabaseMaxConns int32
	DatabaseMinConns int32

	DBMaxConnLifetime time.Duration
	DBConnectTimeout  time.Duration
	DBQueryTimeout    time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{}

	var err error

	if cfg.HTTPAddr, err = getRequiredEnv("HTTP_ADDR"); err != nil {
		return nil, err
	}

	if cfg.DatabaseURL, err = getRequiredEnv("DATABASE_URL"); err != nil {
		return nil, err
	}

	if cfg.LogLevel, err = getRequiredEnv("LOG_LEVEL"); err != nil {
		return nil, err
	}

	if cfg.ShutdownTimeout, err = getDurationEnv("SHUTDOWN_TIMEOUT"); err != nil {
		return nil, err
	}

	if cfg.DatabaseMaxConns, err = getInt32Env("DATABASE_MAX_CONNS"); err != nil {
		return nil, err
	}

	if cfg.DatabaseMinConns, err = getInt32Env("DATABASE_MIN_CONNS"); err != nil {
		return nil, err
	}

	if cfg.DBMaxConnLifetime, err = getDurationEnv("DATABASE_MAX_CONN_LIFETIME"); err != nil {
		return nil, err
	}

	if cfg.DBConnectTimeout, err = getDurationEnv("DATABASE_CONNECT_TIMEOUT"); err != nil {
		return nil, err
	}

	if cfg.DBQueryTimeout, err = getDurationEnv("DATABASE_QUERY_TIMEOUT"); err != nil {
		return nil, err
	}

	return cfg, nil

}

func getRequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}

	return val, nil

}

func getDurationEnv(key string) (time.Duration, error) {
	val, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format for %s: %w", key, err)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return duration, nil
}

func getInt32Env(key string) (int32, error) {
	val, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}

	intVal, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s cannot be converted to int32 format: %w", key, err)
	}

	return int32(intVal), nil
}
