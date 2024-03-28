package models

import "github.com/google/uuid"

type Report struct {
	ID          int       `json:"id"`
	AuthorID    uuid.UUID `json:"authorID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
}
