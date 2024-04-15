package task

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func CreateTask(userID uuid.UUID, task models.Task) (*models.Task, error) {
	var categoryID uuid.UUID

	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO tasks (author_id, title, task_description, points, is_active, category_id) SELECT $1, $2, $3, $4, $5, c.id FROM categories c WHERE c.name = $6 returning id, author_id, title, task_description, category_id, points, is_active",
		userID, task.Title, task.Description, task.Points, task.IsActive, task.Category).Scan(&task.ID, &userID,
		&task.Title, &task.Description, &categoryID, &task.Points, &task.IsActive)
	if err != nil {

		return nil, fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	var categoryName string
	err = db.Pool.QueryRow(context.Background(), "SELECT name FROM categories WHERE id = $1", categoryID).Scan(&categoryName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении названия категории: %w", err)
	}

	task.Category = categoryName

	return &models.Task{
		ID:          task.ID,
		AuthorID:    userID,
		Category:    task.Category,
		Points:      task.Points,
		Title:       task.Title,
		Description: task.Description,
		IsActive:    task.IsActive}, nil
}

func GetTasks() ([]models.Task, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT t.id, title, task_description, is_active, author_id, c.name, points FROM tasks t LEFT JOIN categories c on t.category_id = c.id  WHERE is_active = true")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsActive, &t.AuthorID, &t.Category, &t.Points); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании команды: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
