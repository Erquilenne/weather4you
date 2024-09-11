package middleware

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"

	"github.com/AleksK1NG/api-mc/pkg/utils"
)

func (mw *MiddlewareManager) RequestLoggerMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		rw := negroni.NewResponseWriter(w)

		next(rw, r)

		status := rw.Status()
		size := rw.Size()
		s := time.Since(start).String()
		requestID := utils.GetRequestID(r.Context())

		mw.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s",
			requestID, r.Method, r.URL, status, size, s,
		)
	})
}
