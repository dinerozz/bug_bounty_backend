package models

import "github.com/google/uuid"

type Team struct {
	ID          *int       `json:"id"`
	Name        *string    `json:"name"`
	OwnerID     *uuid.UUID `json:"owner_id"`
	InviteToken *string    `json:"invite_token"`
}

type Teams struct {
	ID   int
	Name string
}

type TeamMember struct {
	UserID      uuid.UUID `json:"id"`
	InviteToken string    `json:"invite_token"`
}
