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
		"INSERT INTO reports (author_id, category_id, title, description) SELECT $1, c.id, $3, $4 FROM categories c WHERE c.name = $2 RETURNING id, author_id, category_id, title, description, status",
		report.AuthorID, report.Category, report.Title, report.Description).Scan(&report.ID, &report.AuthorID, &categoryID, &report.Title, &report.Description, &report.Status)
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

func GetReports(authorID uuid.UUID) ([]models.GetReports, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT r.id, c.name, r.title, r.status from reports r LEFT JOIN categories c on r.category_id = c.id WHERE r.author_id = $1", authorID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отчетов: %w", err)
	}
	defer rows.Close()

	var reports []models.GetReports

	for rows.Next() {
		var r models.GetReports
		if err = rows.Scan(&r.ID, &r.Category, &r.Title, &r.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании отчета")
		}
		reports = append(reports, r)
	}

	return reports, nil
}

func ReviewReport(review models.ReportReview) (*models.ReportReview, error) {
	commandTag, err := db.Pool.Exec(context.Background(), "UPDATE reports SET status = $2 where id = $1", review.ReportID, review.Status)
	if err != nil {
		return nil, fmt.Errorf("ошибка при принятии отчета: %w", err)
	}
	if commandTag.RowsAffected() < 1 {
		return nil, fmt.Errorf("ошибка при принятии отчета: no rows affected")
	}

	row := db.Pool.QueryRow(context.Background(), "INSERT INTO report_reviews (report_id, reviewer_id, review_text, status) VALUES ($1, $2, $3, $4) returning reviewer_id, review_text, status", review.ReportID, review.ReviewerID, review.ReviewText, review.Status)
	if err = row.Scan(&review.ReviewerID, &review.ReviewText, &review.Status); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении вердикта: %w", err)
	}

	return &review, nil
}

func ReviewDetails(reportID int) (*models.ReviewDetails, error) {
	var reviewDetails models.ReviewDetails

	err := db.Pool.QueryRow(context.Background(), "SELECT r.reviewer_id, u.username, r.review_text FROM report_reviews r LEFT JOIN users u on r.reviewer_id = u.id WHERE report_id = $1", reportID).Scan(&reviewDetails.ReviewerID, &reviewDetails.ReviewerUsername, &reviewDetails.ReviewText)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении детального вердикта: %w", err)
	}

	reviewDetails.ReportID = reportID

	return &reviewDetails, nil
}
