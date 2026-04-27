package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/services/go/ewb-gateway-service/internal/api"
	"github.com/complai/complai/services/go/ewb-gateway-service/internal/provider"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	mode := envOr("EWB_PROVIDER", "mock")
	var p provider.EWBProvider
	switch mode {
	case "mock":
		p = provider.NewMockProvider()
		log.Info().Msg("using mock EWB provider")
	default:
		log.Fatal().Str("provider", mode).Msg("unknown EWB provider")
	}

	router := api.NewRouter(p)

	port := envOr("SERVICE_PORT", "8096")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", port).Msg("ewb-gateway-service starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
