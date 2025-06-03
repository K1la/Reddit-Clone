package middleware

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func AccessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" && r.URL.Path != "/manifest.json" && !strings.Contains(r.URL.Path, "/static/") {
			logger.Infow("access log middleware")
			start := time.Now()
			logger.Infow("New request",
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
				"url", r.URL.Path,
				"time", time.Since(start),
			)
		}

		next.ServeHTTP(w, r)

	})
}
