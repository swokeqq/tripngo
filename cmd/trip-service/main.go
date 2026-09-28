package main

import (
	"fmt"
	"log"

	"github.com/swokeqq/tripngo.git/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Println("HTTP Address:", cfg.HTTPAddr)
	fmt.Println("Log Level:", cfg.LogLevel)
	fmt.Println("Shutdown Timeout:", cfg.ShutdownTimeout)
	fmt.Println("Database URL:", cfg.DatabaseURL)
	fmt.Println("Max Conns:", cfg.DatabaseMaxConns)
	fmt.Println("Min Conns:", cfg.DatabaseMinConns)
	fmt.Println("Max Conn Lifetime:", cfg.DBMaxConnLifetime)
	fmt.Println("Connect Timeout:", cfg.DBConnectTimeout)
	fmt.Println("Query Timeout:", cfg.DBQueryTimeout)
}
