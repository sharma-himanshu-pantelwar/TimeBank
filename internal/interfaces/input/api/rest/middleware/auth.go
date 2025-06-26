package middleware

import (
	"context"
	"net/http"
	"timebank/pkg/generatejwt"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("at")
		if err != nil {
			http.Error(w, "missing auth token", http.StatusUnauthorized)
			return
		}

		claims, err := generatejwt.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// set claims in context or request

		ctx := context.WithValue(r.Context(), "user", claims.Uid)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
