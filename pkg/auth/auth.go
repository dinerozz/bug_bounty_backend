package auth

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func RegisterUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	newUserUUID := uuid.New()
	_, err = db.Pool.Exec(context.Background(), "INSERT into users (id, username, email, password) VALUES ($1, $2, $3, $4)", newUserUUID, username, email, hashedPassword)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя в базу данных: %w", err)
	}
	return nil
}

func AuthenticateUser(username, password string) (*models.AuthResponse, error) {
	var (
		userID         uuid.UUID
		hashedPassword string
		email          string
	)

	log.Printf("Попытка аутентификации пользователя: %s", username)

	err := db.Pool.QueryRow(context.Background(), "SELECT id, username, email, password FROM users WHERE username = $1", username).Scan(&userID, &username, &email, &hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("неверный пароль: %w", err)
	}

	accessTTL := time.Now().Add(15 * time.Minute)
	refreshTTL := time.Now().Add(24 * 7 * time.Hour)

	accessToken, refreshToken, err := GenerateTokens(userID)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
		Username:     username,
		Email:        email,
		UserID:       userID,
	}, nil
}

func Refresh(refreshToken string) (*models.AuthResponse, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil || !token.Valid || time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) <= 0 {
		return nil, fmt.Errorf("неверный refresh токен: %w", err)
	}

	accessTokenStr, refreshTokenStr, err := GenerateTokens(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена: %w", err)
	}

	accessTTL := time.Now().Add(15 * time.Minute)
	refreshTTL := time.Now().Add(24 * 7 * time.Hour)

	return &models.AuthResponse{
		Token:        accessTokenStr,
		RefreshToken: refreshTokenStr,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
	}, nil
}

func GenerateTokens(userID uuid.UUID) (accessTokenStr, refreshTokenStr string, err error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	accessExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(24 * time.Hour * 7)
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}
