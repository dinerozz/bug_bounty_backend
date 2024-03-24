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
	}

	authResponse, err := AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ошибка при аутентификации: " + err.Error()})
	}

	accessExpiresInSeconds := int(authResponse.AccessTTL.Sub(time.Now()).Seconds())
	refreshExpiresInSeconds := int(authResponse.RefreshTTL.Sub(time.Now()).Seconds())

	c.SetCookie("auth_token", authResponse.Token, accessExpiresInSeconds, "/", "", true, true)
	c.SetCookie("refresh_token", authResponse.RefreshToken, refreshExpiresInSeconds, "/", "", true, true)

	response := models.CurrentUser{
		ID:       authResponse.UserID,
		Username: authResponse.Username,
		Email:    authResponse.Email,
		Team:     &authResponse.Team,
	}

	c.JSON(http.StatusOK, response)
}

func RefreshHandler(c *gin.Context) {
	cToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token cookie required"})
		return
	}

	tokenString := cToken

	authResponse, err := Refresh(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessExpiresInSeconds := int(authResponse.AccessTTL.Sub(time.Now()).Seconds())
	refreshExpiresInSeconds := int(authResponse.RefreshTTL.Sub(time.Now()).Seconds())

	c.SetCookie("auth_token", authResponse.Token, accessExpiresInSeconds, "/", "", true, true)
	c.SetCookie("refresh_token", authResponse.RefreshToken, refreshExpiresInSeconds, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Токен обновлен"})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Вы успешно вышли из системы"})
}

func CurrentUserHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "UserID type assertion failed"})
		return
	}

	user, err := GetUserByID(db.Pool, userIDStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
