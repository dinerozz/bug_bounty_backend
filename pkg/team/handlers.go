package team

import (
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func CreateTeamHandler(c *gin.Context) {
	var newTeam models.Team

	if err := c.ShouldBindJSON(&newTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find user ID"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID format is incorrect"})
		return
	}

	newTeam.OwnerID = userID

	err := CreateTeam(&newTeam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTeam)
}

func GetTeamsHandler(c *gin.Context) {
	teams, err := GetTeams()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Произошла ошибка при получении списка команд: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, teams)
}
