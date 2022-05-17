package middlewares

import (
	"context"
	"net/http"

	"github.com/csothen/tmdei-project/internal/services"
)

func Auth(s *services.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if auth == "" {
				next.ServeHTTP(w, r)
				return
			}

			err := s.ValidateAuth(context.Background(), auth)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
