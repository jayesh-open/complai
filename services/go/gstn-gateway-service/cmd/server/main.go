package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/api"
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/provider"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	providerMode := envOr("GSTN_PROVIDER", "mock")

	var p provider.GSTNProvider
	switch providerMode {
	case "mock":
		log.Info().Msg("using mock GSTN provider")
		p = provider.NewMockProvider()
	case "adaequare":
		log.Fatal().Msg("adaequare provider not yet implemented — set GSTN_PROVIDER=mock")
	default:
		log.Fatal().Str("provider", providerMode).Msg("unknown GSTN_PROVIDER")
	}

	router := api.NewRouter(p)

	port := envOr("SERVICE_PORT", "8091")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", port).Str("provider", providerMode).Msg("gstn-gateway-service starting")
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

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
