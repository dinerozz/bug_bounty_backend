package team

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
)

func CreateTeam(team *models.Team) error {
	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO teams (name, owner_id) VALUES ($1, $2) RETURNING id",
		team.Name, team.OwnerID).Scan(&team.ID)

	if err != nil {
		return fmt.Errorf("ошибка при создании команды: %w", err)
	}

	return nil
}

func GetTeams() ([]models.Teams, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, name FROM teams")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении команд: %w", err)
	}
	defer rows.Close()

	var teams []models.Teams

	for rows.Next() {
		var t models.Teams
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании команды: %w", err)
		}
		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после итерации по командам: %w", err)
	}

	return teams, nil
}
