package api

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("health-probe-service")

var startTime = time.Now()

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	GoVersion string `json:"go_version"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := tracer.Start(ctx, "health-check", trace.WithAttributes(
		attribute.String("service", "health-probe-service"),
	))
	defer span.End()

	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "0.1.0-dev"
	}

	resp := HealthResponse{
		Status:    "healthy",
		Version:   version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		GoVersion: runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
