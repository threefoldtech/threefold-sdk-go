package prometheus_integration

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Metrics struct {
	RequestCount   *prometheus.CounterVec
	RequestLatency prometheus.Histogram
	SuccessRate    prometheus.Gauge
	TotalSuccess   prometheus.Counter
	ErrorRate      prometheus.Gauge
}

// NewMetrics returns a new Metrics instance that must be
// registered in a Prometheus registry with Register.
func NewMetrics() *Metrics {
	return &Metrics{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests categorized by status",
			},
			[]string{"status"},
		),
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
		TotalSuccess: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_requests_success_total",
				Help: "Total number of successful HTTP requests",
			},
		),
		ErrorRate: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_error_rate",
				Help: "Percentage of failed HTTP requests",
			},
		),
	}
}

// Register registers all Prometheus metrics in r.
func (pm *Metrics) Register(r prometheus.Registerer) error {
	for _, c := range []prometheus.Collector{
		pm.TotalSuccess,
		pm.SuccessRate,
		pm.RequestCount,
		pm.ErrorRate,
		pm.RequestLatency,
	} {
		if err := r.Register(c); err != nil {
			return fmt.Errorf("failed to register metric %v: %w", c, err)
		}
	}
	return nil
}

// Observe metrics given a vegeta.Result.
func (pm *Metrics) Observe(res *vegeta.Result) {
	pm.RequestCount.WithLabelValues(fmt.Sprintf("%d", res.Code)).Inc()
	pm.RequestLatency.Observe(res.Latency.Seconds())

	if res.Code >= 200 && res.Code < 300 {
		pm.TotalSuccess.Inc()
	}

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
