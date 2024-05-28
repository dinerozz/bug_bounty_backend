package models

import "github.com/google/uuid"

type Conversation struct {
	ID       uuid.UUID `json:"id"`
	ReportID int       `json:"report_id"`
	UserID   uuid.UUID `json:"user_id"`
	Message  string    `json:"message"`
}

type GetConversation struct {
	ID       uuid.UUID `json:"id"`
	ReportID int       `json:"report_id"`
	UserID   uuid.UUID `json:"user_id"`
	Message  string    `json:"message"`
	Username string    `json:"username"`
	IsAdmin  bool      `json:"is_admin"`
}
