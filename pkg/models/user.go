package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Points   *int      `json:"points"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type CurrentUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Points   *int      `json:"points"`
	Email    string    `json:"email"`
	Team     *Team     `json:"team,omitempty"`
}
