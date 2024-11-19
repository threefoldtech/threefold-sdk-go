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
	SuccessForEachEndpoint      *prometheus.CounterVec
	FailForEachEndpoint         *prometheus.CounterVec
	TotalRequestForEachEndpoint *prometheus.CounterVec
	RequestLatency              prometheus.Histogram
	SuccessRate                 prometheus.Gauge
	ErrorRate                   prometheus.Gauge
	AvgLatency                  prometheus.Gauge
	MaxLatency                  prometheus.Gauge
}

// NewMetrics returns a new Metrics instance that must be
// registered in a Prometheus registry with Register.
func NewMetrics() *Metrics {
	return &Metrics{
		SuccessForEachEndpoint: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "successes_total",
				Help: "Total successful requests for each endpoint.",
			},
			[]string{"endpoint"},
		),
		FailForEachEndpoint: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "failures_total",
				Help: "Total failed requests for each endpoint.",
			},
			[]string{"endpoint"},
		),

		TotalRequestForEachEndpoint: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_requests",
				Help: "Total requests made to each endpoint.",
			},
			[]string{"endpoint"},
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
	}
}

// Register registers all Prometheus metrics in r.
func (pm *Metrics) Register(r prometheus.Registerer) error {
	for _, c := range []prometheus.Collector{
		pm.SuccessForEachEndpoint,
		pm.FailForEachEndpoint,
		pm.TotalRequestForEachEndpoint,
		pm.SuccessRate,
		pm.ErrorRate,
		pm.RequestLatency,
		pm.AvgLatency,
		pm.MaxLatency,
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

	if res.Code >= 200 && res.Code < 300 {
		pm.SuccessForEachEndpoint.WithLabelValues(res.URL).Inc()
	} else {
		pm.FailForEachEndpoint.WithLabelValues(res.URL).Inc()
	}

	pm.TotalRequestForEachEndpoint.WithLabelValues(res.URL).Inc()

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
