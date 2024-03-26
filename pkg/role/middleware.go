package role

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func RolesMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}

		userID, _ := userIDInterface.(uuid.UUID)

		userRole, err := GetUserRole(userID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
