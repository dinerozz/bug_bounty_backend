package models

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type UserRole struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}
