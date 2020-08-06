package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Timer(m prometheus.Observer, next http.HandlerFunc) http.HandlerFunc {
	if m == nil {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		go m.Observe(time.Since(start).Seconds())
	}
}

func Counter(m prometheus.Counter, next http.HandlerFunc) http.HandlerFunc {
	if m == nil {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		go m.Inc()
		next(w, r)
	}
}
