package middleware

import (
	"net/http"

	"github.com/simoncdn/http-server/internal/config"
)

func Metrics(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.IncrementHits()
			next.ServeHTTP(w, r)
		})
	}
}
