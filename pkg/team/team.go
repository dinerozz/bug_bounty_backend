package team

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func CreateTeam(team *models.Team) error {
	inviteToken, _ := generateRandomString(32)
	fmt.Println(inviteToken)

	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO teams (name, owner_id, invite_token) VALUES ($1, $2, $3) RETURNING id, invite_token",
		team.Name, team.OwnerID, inviteToken).Scan(&team.ID, &inviteToken)

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

func GetTeamByOwner(userID uuid.UUID) (models.Team, error) {
	var team models.Team

	err := db.Pool.QueryRow(context.Background(), "SELECT * FROM teams WHERE owner_id = $1", userID).Scan(
		&team.ID,
		&team.Name,
		&team.OwnerID,
	)

	if err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func generateRandomString(length int) (string, error) {
	bytesLength := (length * 3) / 4

	randomBytes := make([]byte, bytesLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(randomBytes), nil
}
