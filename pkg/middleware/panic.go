package middleware

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func Panic(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" && r.URL.Path != "/manifest.json" && !strings.Contains(r.URL.Path, "/static/") {
			logger.Infow("panicMiddleware", r.URL.Path)
		}
		defer func() {
			if err := recover(); err != nil {
				logger.Infow("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
