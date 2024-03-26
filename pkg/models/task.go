package models

import "github.com/google/uuid"

type Task struct {
	ID          int       `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Points      int       `json:"points"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
}
