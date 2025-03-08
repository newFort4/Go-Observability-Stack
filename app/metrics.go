package main

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
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
