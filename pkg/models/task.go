package models

import "github.com/google/uuid"

type Task struct {
	ID          int       `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
}
