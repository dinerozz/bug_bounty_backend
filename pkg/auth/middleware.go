package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "Authorization cookie required", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenString := c.Value

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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
