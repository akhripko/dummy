package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HelloRequestCounts prometheus.Counter
	HelloRequestTiming prometheus.Histogram
	SomeGauge          prometheus.Gauge
)

func Register() {
	HelloRequestCounts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hello_requests_count",
		Help: "hello requests count",
	})

	HelloRequestTiming = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "hello_request_timing",
		Help:    "hello request timing",
		Buckets: prometheus.ExponentialBuckets(0.5, 2, 15),
	})

	SomeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "some_gauge",
		Help: "The gauge for something",
	})

	prometheus.MustRegister(
		HelloRequestTiming,
		HelloRequestCounts,
		SomeGauge,
	)
}

func Unregister() { // nolint megacheck
	prometheus.Unregister(HelloRequestTiming)
	prometheus.Unregister(HelloRequestCounts)
	prometheus.Unregister(SomeGauge)
}
