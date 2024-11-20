package prometheus_integration

import (
	"fmt"
	"net/http"
	"time"

	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	vegeta "github.com/tsenart/vegeta/lib"
)

type Metrics struct {
	RequestLatency prometheus.Histogram
	SuccessRate    prometheus.Gauge
	ErrorRate      prometheus.Gauge
	AvgLatency     prometheus.Gauge
	MaxLatency     prometheus.Gauge
	RequestCode    *prometheus.CounterVec
}

// NewMetrics returns a new Metrics instance that must be
// registered in a Prometheus registry with Register.
func NewMetrics() *Metrics {
	return &Metrics{
		RequestLatency: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_latency_seconds",
				Help:    "Histogram of response latencies",
				Buckets: prometheus.DefBuckets,
			},
		),
		SuccessRate: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_success_rate",
				Help: "Percentage of successful HTTP requests",
			},
		),
		ErrorRate: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_error_rate",
				Help: "Percentage of failed HTTP requests",
			},
		),
		AvgLatency: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "avg_latency",
			},
		),
		MaxLatency: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "max_latency",
			},
		),
		RequestCode: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_http_request_with_status_response",
			}, []string{"status"},
		),
	}
}

// Register registers all Prometheus metrics in r.
func (pm *Metrics) Register(r prometheus.Registerer) error {
	for _, c := range []prometheus.Collector{
		pm.SuccessRate,
		pm.ErrorRate,
		pm.RequestLatency,
		pm.AvgLatency,
		pm.MaxLatency,
		pm.RequestCode,
	} {
		if err := r.Register(c); err != nil {
			return fmt.Errorf("failed to register metric %v: %w", c, err)
		}
	}
	return nil
}

// Observe metrics given a vegeta.Result.
func (pm *Metrics) Observe(res *vegeta.Result) {
	pm.RequestLatency.Observe(res.Latency.Seconds())

	pm.RequestCode.WithLabelValues(strconv.Itoa(int(res.Code))).Inc()

}

// NewHandler returns a new http.Handler that exposes Prometheus
// metrics registered in r in the OpenMetrics format.
func NewHandler(r *prometheus.Registry, startTime time.Time) http.Handler {
	return promhttp.HandlerFor(r, promhttp.HandlerOpts{
		Registry:          r,
		EnableOpenMetrics: true,
		ProcessStartTime:  startTime,
	})
}
