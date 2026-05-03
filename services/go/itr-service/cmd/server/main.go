package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/complai/complai/packages/shared-kernel-go/db"
	"github.com/complai/complai/services/go/itr-service/internal/api"
	pgstore "github.com/complai/complai/services/go/itr-service/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	port := envOr("SERVICE_PORT", "8100")

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	pool, err := db.NewPool(context.Background(), envOr("DATABASE_URL", "postgres://complai:complai_dev@localhost:5432/itr_db?sslmode=disable"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	repo := pgstore.NewPgStore(pool)
	h := api.NewHandlers(repo)
	router := api.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", port).Msg("itr-service listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("shutdown failed")
	}
	log.Info().Msg("itr-service stopped")
}

func runMigrations() {
	dbURL := envOr("DATABASE_URL", "postgres://complai:complai_dev@localhost:5432/itr_db?sslmode=disable")
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db for migrations")
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("failed to set goose dialect")
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}
	log.Info().Msg("itr-service migrations complete")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
