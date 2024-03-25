package role

import (
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRoleHandler(c *gin.Context) {
	var request models.Role
	var role *models.Role

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := CreateRole(models.Role{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, role)
}

func SetUserRoleHandler(c *gin.Context) {
	var request models.UserRole

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := SetUserRole(models.UserRole{
		UserID: request.UserID,
		Role:   request.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось установить роль"})
		return
	}

	c.JSON(http.StatusOK, request)
}
