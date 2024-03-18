package models

import "time"

type Invite struct {
	ID        int
	TeamID    int
	Token     string
	ExpiresAt time.Time
}
