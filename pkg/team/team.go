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
