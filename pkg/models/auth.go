package models

import (
	"github.com/google/uuid"
	"time"
)

type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	AccessTTL    time.Time `json:"expiresAt"`
	RefreshTTL   time.Time `json:"expiresAt"`
	UserID       uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
}

type RegisterBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthBody struct {
	Username string `json: "username"`
	Password string `json:"password"`
}
