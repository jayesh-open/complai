package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/complai/complai/packages/shared-kernel-go/db"
	"github.com/complai/complai/services/go/workflow-service/internal/api"
	"github.com/complai/complai/services/go/workflow-service/internal/store"
	"github.com/complai/complai/services/go/workflow-service/internal/workflows"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	ctx := context.Background()

	dbURL := envOr("DATABASE_URL", "postgres://complai_app:complai_app_dev@localhost:5432/workflow_db?sslmode=disable")
	pool, err := db.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("db pool failed")
	}
	defer pool.Close()

	s := store.New(pool)

	// Temporal client
	temporalAddr := envOr("TEMPORAL_ADDRESS", "localhost:7233")
	temporalNS := envOr("TEMPORAL_NAMESPACE", "default")

	var engine api.WorkflowEngine
	var tc client.Client
	tc, err = client.Dial(client.Options{
		HostPort:  temporalAddr,
		Namespace: temporalNS,
	})
	if err != nil {
		log.Warn().Err(err).Msg("temporal connect failed, workflows disabled")
		engine = api.NewNoopEngine()
	} else {
		defer tc.Close()
		engine = api.NewTemporalEngine(tc)
	}

	router := api.NewRouter(s, engine)

	// Start Temporal worker
	if tc != nil {
		w := worker.New(tc, "complai-default", worker.Options{})
		w.RegisterWorkflow(workflows.SampleSagaWorkflow)
		w.RegisterActivity(workflows.ValidateActivity)
		w.RegisterActivity(workflows.CreateHumanTaskActivity)
		w.RegisterActivity(workflows.FinalizeActivity)
		go func() {
			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Error().Err(err).Msg("temporal worker failed")
			}
		}()
	}

	port := envOr("SERVICE_PORT", "8089")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", port).Msg("workflow-service starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Fatal().Err(err).Msg("forced shutdown")
	}
}

func runMigrations() {
	dbURL := envOr("MIGRATION_DATABASE_URL", envOr("DATABASE_URL", "postgres://complai:complai_dev@localhost:5432/workflow_db?sslmode=disable"))
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("open db for migration")
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("set dialect")
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("run migrations")
	}
	log.Info().Msg("workflow-service migrations complete")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
