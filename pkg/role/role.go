package role

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func CreateRole(role models.Role) (*models.Role, error) {
	var newRole models.Role

	err := db.Pool.QueryRow(context.Background(), "INSERT INTO roles (name, description) VALUES ($1, $2) returning id, name, description", role.Name, role.Description).Scan(&newRole.ID, &newRole.Name, &newRole.Description)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании роли: %w", err)

	}

	return &newRole, nil
}

func SetUserRole(request models.UserRole) error {
	_, err := db.Pool.Exec(context.Background(), "INSERT INTO user_roles (user_id, role_id) SELECT $1, id from roles WHERE roles.name = $2", request.UserID, request.Role)
	if err != nil {
		fmt.Println("error", err)
		return fmt.Errorf("произошла ошибка при выдаче роли: %w", err)
	}

	return nil
}

func GetUserRole(userID uuid.UUID) (string, error) {
	var userRole string

	err := db.Pool.QueryRow(context.Background(),
		"SELECT r.name FROM roles r LEFT JOIN user_roles ur on r.id = ur.role_id WHERE ur.user_id = $1", userID).Scan(&userRole)

	if err != nil {
		return "", fmt.Errorf("не удалось получить роль пользователя: %w", err)

	}

	return userRole, nil
}
