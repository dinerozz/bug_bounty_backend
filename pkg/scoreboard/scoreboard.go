package scoreboard

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
)

func GetScoreboard() ([]models.ScoreBoard, error) {
	var scoreboard []models.ScoreBoard

	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, points FROM teams")
	if err != nil {
		fmt.Errorf("ошибка при получении записей по командам: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ScoreBoard
		if err = rows.Scan(&s.ID, &s.Name, &s.Points); err != nil {
			return nil, err
		}
		scoreboard = append(scoreboard, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return scoreboard, nil
}
