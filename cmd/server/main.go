package main

import (
	"cm_open_api/internal/config"
	"cm_open_api/internal/handlers"
	"cm_open_api/internal/metrics"
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize shared pgx connection pool
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, cfg.PostgresConnStr)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer pool.Close()
	cfg.DB = pool

	go metrics.SetupPrometheus()

	r := handlers.SetupRouter(cfg)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
