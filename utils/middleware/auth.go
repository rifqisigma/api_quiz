package middleware

import (
	"api_quiz/utils/helper"
	"context"
	"net/http"
	"strings"
)

type key int

const UserContextKey key = 0

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.ParseJWT(tokenString)
		if err != nil || !claims.IsVerified {
			http.Error(w, "Unauthorized: Invalid or unverified user", http.StatusForbidden)
			return
		}

		if !claims.IsVerified {
			http.Error(w, "Unauthorized: user is not verified", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
