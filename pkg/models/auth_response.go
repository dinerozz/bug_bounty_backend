package models

import "time"

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
