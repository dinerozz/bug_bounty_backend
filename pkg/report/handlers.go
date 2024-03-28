package report

import (
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func CreateReportHandler(c *gin.Context) {
	userIDInterface, _ := c.Get("userID")
	var request models.Report

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "UserID type assertion failed"})
		return
	}

	request.AuthorID = userID

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := CreateReport(request)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании отчета"})
	}

	c.JSON(http.StatusOK, report)
}

func ReviewReportHandler(c *gin.Context) {
	userIDInterface, _ := c.Get("userID")

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "UserID type assertion failed"})
		return
	}

	var request models.ReportReview
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.ReviewerID = userID

	review, err := ReviewReport(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при публикации вердикта"})
	}

	c.JSON(http.StatusOK, review)
}