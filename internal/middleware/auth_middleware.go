package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shariski/room-booking/internal/model"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				auth = auth[7:]
			}

			token, err := jwt.ParseWithClaims(auth, &model.Auth{}, func(t *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*model.Auth)
			if !ok || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
