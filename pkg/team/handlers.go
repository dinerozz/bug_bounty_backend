package team

import (
	"fmt"
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

	newTeam.OwnerID = &userID

	team, err := CreateTeam(&newTeam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(*team.Name) < 2 {
		fmt.Println(team.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Слишком короткое название команды"})
		return
	}

	c.JSON(http.StatusCreated, team)
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

func UpdateInviteTokenHandler(c *gin.Context) {
	userIDInterface, ok := c.Get("userID")

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "UserID type assertion failed"})
		return
	}

	inviteToken, err := UpdateInviteToken(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "произошла ошибка при обновлении токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite_token": inviteToken})
}
