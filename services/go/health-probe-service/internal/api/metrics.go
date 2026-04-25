package api

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "# HELP health_probe_up Whether the service is up\n")
	fmt.Fprintf(w, "# TYPE health_probe_up gauge\n")
	fmt.Fprintf(w, "health_probe_up 1\n\n")

	fmt.Fprintf(w, "# HELP health_probe_uptime_seconds Service uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE health_probe_uptime_seconds gauge\n")
	fmt.Fprintf(w, "health_probe_uptime_seconds %.0f\n\n", time.Since(startTime).Seconds())

	fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use\n")
	fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n\n", m.Alloc)

	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
}
