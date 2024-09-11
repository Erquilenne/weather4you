package middleware

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"

	"weather4you/pkg/metric"
)

// Prometheus metrics middleware
func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		var status int
		metrics.ObserveResponseTime(status, r.Method, r.URL.Path, time.Since(start).Seconds())
		metrics.IncHits(status, r.Method, r.URL.Path)
	})
}
