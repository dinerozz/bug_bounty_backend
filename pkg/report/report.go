package report

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func CreateReport(report models.Report) (*models.Report, error) {
	var categoryID uuid.UUID

	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO reports (author_id, category_id, title, description) SELECT $1, c.id, $3, $4 FROM categories c WHERE c.name = $2 RETURNING id, author_id, category_id, title, description",
		report.AuthorID, report.Category, report.Title, report.Description).Scan(&report.ID, &report.AuthorID, &categoryID, &report.Title, &report.Description)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании отчета: %w", err)
	}

	var categoryName string
	err = db.Pool.QueryRow(context.Background(), "SELECT name FROM categories WHERE id = $1", categoryID).Scan(&categoryName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении названия категории: %w", err)
	}

	report.Category = categoryName

	return &report, nil
}
