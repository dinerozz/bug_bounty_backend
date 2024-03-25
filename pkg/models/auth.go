package models

import (
	"time"
)

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	AccessTTL    time.Time    `json:"expiresAt"`
	RefreshTTL   time.Time    `json:"expiresAt"`
	User         *CurrentUser `json:"current_user"`
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
