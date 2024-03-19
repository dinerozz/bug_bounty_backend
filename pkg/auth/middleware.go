package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(authHeaderParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
