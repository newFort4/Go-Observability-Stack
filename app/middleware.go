package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func defaultMiddleware(httpHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a tracer
		tracer := otel.Tracer("example-tracer")

		// Start a trace span
		ctx, span := tracer.Start(context.Background(), "example-span")
		span.SetAttributes(attribute.String("key", "value")) // Add metadata
		defer span.End()

		time.Sleep(100 * time.Millisecond)

		_, child := tracer.Start(ctx, "child-span")
		span.SetAttributes(attribute.String("key", "value")) // Add metadata
		defer child.End()

		requests.WithLabelValues(r.Method).Inc()
		start := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method))

		logger.Info("Request received", "method", r.Method, "path", r.URL.Path)

		httpHandler(w, r)

		logger.Info("Request completed", "method", r.Method, "path", r.URL.Path, slog.Int("statusCode", r.Response.StatusCode))

		if r.Response.StatusCode >= 400 {
			errorCounter.WithLabelValues(r.Method, fmt.Sprintf("%d", r.Response.StatusCode)).Inc()
		} else {
			requestSize.WithLabelValues(r.Method).Observe(float64(r.Response.ContentLength))
		}
		start.ObserveDuration()
	}
}
