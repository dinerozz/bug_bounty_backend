package auth

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var jwtKey = []byte("8149694d8d0bfcdddc2c965b2a2f2ba1d4233eb778901ca7882651f291eb828a")

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

func AuthenticateUser(username, password string) (string, error) {
	var userID uuid.UUID
	var hashedPassword string

	log.Printf("Попытка аутентификации пользователя: %s", username)

	err := db.Pool.QueryRow(context.Background(), "SELECT id, password FROM users WHERE username = $1", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("пользователь не найден: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("неверный пароль: %w", err)
	}

	expirationTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", fmt.Errorf("ошибка при создании токена: %w", err)
	}

	return tokenString, nil
}
