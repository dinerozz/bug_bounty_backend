package models

import "github.com/google/uuid"

type Report struct {
	ID          int       `json:"id"`
	AuthorID    uuid.UUID `json:"authorID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
}

type ReportReview struct {
	ReportID   int       `json:"report_id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	ReviewText string    `json:"review_text"`
	Status     string    `json:"status"`
}

type ReviewDetails struct {
	ReportID         int       `json:"report_id"`
	ReviewerID       uuid.UUID `json:"reviewer_id"`
	ReviewerUsername string    `json:"reviewer_username"`
	ReviewText       string    `json:"review_text"`
}

type GetReports struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Status   string `json:"status"`
}
