package models

import "time"

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
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
