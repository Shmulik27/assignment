package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	// Processing metrics
	ProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "processing_duration_seconds",
			Help:    "Time taken to process each record",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Worker metrics
	ActiveWorkers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name:    "active_workers",
			Help:    "Number of currently active workers",
		},
	)

	// Channel metrics
	ChannelCapacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:    "channel_capacity",
			Help:    "Current capacity of processing channels",
		},
		[]string{"channel"},
	)

	// Error metrics
	ProcessingErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:    "processing_errors_total",
			Help:    "Total number of processing errors",
		},
		[]string{"type"},
	)
)

// TrackDuration measures the duration of an operation
func TrackDuration(operation string, start time.Time) {
	ProcessingDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
}

// RecordError increments the error counter for a specific error type
func RecordError(errorType string) {
	ProcessingErrors.WithLabelValues(errorType).Inc()
}

// UpdateChannelCapacity updates the channel capacity metric
func UpdateChannelCapacity(channelName string, capacity float64) {
	ChannelCapacity.WithLabelValues(channelName).Set(capacity)
} 