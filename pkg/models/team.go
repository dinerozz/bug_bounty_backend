package models

import "github.com/google/uuid"

type Team struct {
	ID      int
	Name    string
	OwnerID uuid.UUID
}
