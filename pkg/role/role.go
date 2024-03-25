package role

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
)

func CreateRole(role models.Role) (*models.Role, error) {
	var newRole models.Role

	err := db.Pool.QueryRow(context.Background(), "INSERT INTO roles (name, description) VALUES ($1, $2) returning id, name, description", role.Name, role.Description).Scan(&newRole.ID, &newRole.Name, &newRole.Description)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании роли: %w", err)

	}

	return &newRole, nil
}
