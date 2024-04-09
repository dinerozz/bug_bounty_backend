package scoreboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetScoreboardHandler(c *gin.Context) {
	scoreboard, err := GetScoreboard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при получении таблицы"})
		return
	}

	c.JSON(http.StatusOK, scoreboard)
}
