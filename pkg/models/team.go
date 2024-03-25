package models

import "github.com/google/uuid"

type Team struct {
	ID          *int       `json:"id"`
	Name        *string    `json:"name"`
	OwnerID     *uuid.UUID `json:"owner_id"`
	InviteToken *string    `json:"invite_token"`
	TeamMembers []Member   `json:"members"`
}

type Teams struct {
	ID   int
	Name string
}

type JoinTeam struct {
	UserID      uuid.UUID `json:"id"`
	InviteToken string    `json:"invite_token"`
}

type Member struct {
	ID       uuid.UUID `json:"id"'`
	Username string    `json:"username"'`
}
