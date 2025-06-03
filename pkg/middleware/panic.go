package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" && r.URL.Path != "/manifest.json" && !strings.Contains(r.URL.Path, "/static/") {
			fmt.Println("panicMiddleware", r.URL.Path)
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
