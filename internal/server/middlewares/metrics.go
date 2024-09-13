package middlewares

import (
	"clean-rest-arch/internal/server/utils/metrics"
	"net/http"
	"time"
)

func MiddlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		metrics.RequestDuration.WithLabelValues(path, method).Observe(float64(duration))
		metrics.RequestsTotal.WithLabelValues(path, method, w.Header().Get("status")).Inc()
	})
}
