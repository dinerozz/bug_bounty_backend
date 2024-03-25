package task

import (
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func CreateTaskHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	var request models.Task

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "UserID type assertion failed"})
		return
	}

	task, err := CreateTask(userID, request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTasksHandler(c *gin.Context) {
	tasks, err := GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении задач"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
