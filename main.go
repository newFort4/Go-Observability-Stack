package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	logger   *slog.Logger
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_request_duration_seconds",
			Help:    "Histogram of request durations for the application",
			Buckets: prometheus.DefBuckets, // Default bucket sizes
		},
		[]string{"method"},
	)
	errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_request_errors_total",
			Help: "Total number of errors encountered by the application",
		},
		[]string{"method", "status"},
	)
	requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_app_response_size_bytes",
			Help:    "Histogram of the sizes of responses",
			Buckets: prometheus.ExponentialBuckets(100, 2, 4), // Example of custom buckets
		},
		[]string{"method"},
	)
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requests)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(requestSize)
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	logger.Info("Request from user", "name", name)

	fmt.Fprintf(w, "Hello, %s!\n", name)
	r.Response = new(http.Response)
	if rand.Int()%2 == 0 {
		r.Response.StatusCode = http.StatusInternalServerError
	} else {
		r.Response.StatusCode = http.StatusOK
	}

	r.Response.ContentLength = int64(len(name) + 8)

	// Generate a random duration between 100ms and 500ms
	sleepDuration := time.Duration(rand.Intn(401)+100) * time.Millisecond
	time.Sleep(sleepDuration)
}

func defaultMiddleware(httpHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a tracer
		tracer := otel.Tracer("example-tracer")

		// Start a trace span
		_, span := tracer.Start(context.Background(), "example-span")
		span.SetAttributes(attribute.String("key", "value")) // Add metadata
		defer span.End()

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

func initTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("tempo:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("go-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

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
