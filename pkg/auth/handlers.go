package auth

import (
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func RegisterHandler(c *gin.Context) {
	var req models.RegisterBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
		return
	}

	if err := RegisterUser(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь " + req.Username + " успешно зарегистрирован"})
}

func AuthenticateHandler(c *gin.Context) {
	var req models.AuthBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
		return
	}

	authResponse, err := AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка при аутентификации: " + err.Error()})
		return
	}

	expiresInSeconds := int(authResponse.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("auth_token", authResponse.Token, expiresInSeconds, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация"})
}

func RefreshHandler(c *gin.Context) {
	cToken, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token cookie required"})
		return
	}

	tokenString := cToken

	authResponse, err := Refresh(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expiresInSeconds := int(authResponse.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("auth_token", authResponse.Token, expiresInSeconds, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Токен обновлен"})
}

func CurrentUserHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID type assertion failed"})
		return
	}

	user, err := db.GetUserByID(db.Pool, userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
