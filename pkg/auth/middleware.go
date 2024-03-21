package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtKey := []byte(os.Getenv("JWT_KEY"))

		cToken, err := c.Cookie("auth_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie required"})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		tokenString := cToken

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
