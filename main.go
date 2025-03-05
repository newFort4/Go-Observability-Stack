package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logFile, err := os.OpenFile("/var/log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger = slog.New(slog.NewJSONHandler(logFile, nil))

	tp, err := initTracer()
	if err != nil {
		log.Fatal("failed to initialize OpenTelemetry:", err)
	}
	defer tp.Shutdown(context.Background())

	http.HandleFunc("/", defaultMiddleware(handler))
	// Expose the /metrics endpoint for Prometheus to scrape
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
