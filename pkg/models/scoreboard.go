package models

type ScoreBoard struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Points *int   `json:"points"`
}
