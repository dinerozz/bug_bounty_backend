package task

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func CreateTask(userID uuid.UUID, task models.Task) (*models.Task, error) {
	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO tasks (author_id, title, task_description, is_active) VALUES ($1, $2, $3, $4) returning id, author_id, title, task_description, is_active",
		userID, task.Title, task.Description, task.IsActive).Scan(&task.ID, &userID,
		&task.Title, &task.Description, &task.IsActive)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании команды: %w", err)
	}

	return &models.Task{
		ID:          task.ID,
		AuthorID:    userID,
		Title:       task.Title,
		Description: task.Description,
		IsActive:    task.IsActive}, nil
}

func GetTasks() ([]models.Task, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, task_description, is_active, author_id FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsActive, &t.AuthorID); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании команды: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
