package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry *prometheus.Registry
)

type Metrics struct {
	HTTPRequestsReceived          *prometheus.CounterVec
	DBOperationsErrors            *prometheus.CounterVec
	InternalErrors                *prometheus.CounterVec
	DBOperationsDuration          *prometheus.HistogramVec
	HTTPRequestProcessingDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_received",
			Help: "Number of HTTP requests received by the server",
		}, []string{"method", "path"}),
		DBOperationsErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "db_operations_errors",
			Help: "Number of errors encountered during database operations",
		}, []string{"operation", "table"}),
		InternalErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "internal_errors",
			Help: "Number of internal errors encountered by the server",
		}, []string{"error"}),
		DBOperationsDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "db_operations_duration",
			Help: "Duration of database operations",
		}, []string{"operation", "table"}),
		HTTPRequestProcessingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_processing_duration",
			Help: "Duration of HTTP request processing",
		}, []string{"method", "path"}),
	}
}

func (m *Metrics) RecordDuration(h *prometheus.HistogramVec, labels []string) func() {
	start := time.Now()
	return func() {
		h.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) RecordCount(c *prometheus.CounterVec, labels []string) {
	c.WithLabelValues(labels...).Inc()
}

func (m *Metrics) Register() error {
	registry = prometheus.NewRegistry()
	for _, metric := range []prometheus.Collector{
		m.HTTPRequestsReceived,
		m.DBOperationsErrors,
		m.InternalErrors,
		m.DBOperationsDuration,
		m.HTTPRequestProcessingDuration,
	} {
		if err := registry.Register(metric); err != nil {
			return fmt.Errorf("failed to register metric: %w", err)
		}
	}
	return nil
}

func (m *Metrics) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry:          registry,
		EnableOpenMetrics: true,
	})
}
