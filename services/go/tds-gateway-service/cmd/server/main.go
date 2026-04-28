package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/api"
	"github.com/complai/complai/services/go/tds-gateway-service/internal/provider"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	mode := envOr("SANDBOX_MODE", "mock")
	port := envOr("SERVICE_PORT", "8098")

	var p provider.SandboxTDSProvider
	switch mode {
	case "mock":
		log.Info().Msg("starting tds-gateway in mock mode")
		p = provider.NewMockProvider()
	case "real":
		log.Fatal().Msg("real Sandbox.co.in TDS provider not yet implemented — waiting for auth resolution")
	default:
		log.Fatal().Str("mode", mode).Msg("unknown SANDBOX_MODE")
	}

	h := api.NewHandlers(p)
	router := api.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", port).Msg("tds-gateway-service listening")
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
	log.Info().Msg("tds-gateway-service stopped")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
