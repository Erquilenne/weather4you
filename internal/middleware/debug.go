package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/urfave/negroni"
)

// Debug dump request middleware for Negroni
func (mw *MiddlewareManager) DebugMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if mw.cfg.Server.Debug {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			mw.logger.Info(fmt.Sprintf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump))
		}
		next(w, r)
	})
}
