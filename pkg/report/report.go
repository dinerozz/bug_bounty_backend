package report

import (
	"context"
	"errors"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strconv"
)

func CreateReport(report models.Report) (*models.Report, error) {
	var categoryID uuid.UUID

	err := db.Pool.QueryRow(context.Background(),
		"INSERT INTO reports (author_id, category_id, title, description, team_id) SELECT $1, c.id, $3, $4, $5 FROM categories c WHERE c.name = $2 RETURNING id, author_id, category_id, title, description, status, team_id",
		report.AuthorID, report.Category, report.Title, report.Description, report.TeamID).Scan(&report.ID, &report.AuthorID, &categoryID, &report.Title, &report.Description, &report.Status, &report.TeamID)
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
	rows, err := db.Pool.Query(context.Background(), ReportsTableQuery, authorID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отчетов: %w", err)
	}
	defer rows.Close()

	var reports []models.GetReports

	var id = 0
	for rows.Next() {
		var r models.GetReports
		id++
		if err = rows.Scan(&r.ReportID, &r.Category, &r.Title, &r.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании отчета")
		}
		r.ID = id
		reports = append(reports, r)
	}

	return reports, nil
}

func GetAdminReports() ([]models.GetReports, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT r.id, c.name, r.title, r.status FROM reports r LEFT JOIN categories c on r.category_id = c.id order by r.id")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отчетов: %w", err)
	}

	var id = 0
	var reports []models.GetReports
	for rows.Next() {
		var r models.GetReports
		id++
		if err = rows.Scan(&r.ReportID, &r.Category, &r.Title, &r.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании отчета")
		}
		r.ID = id
		reports = append(reports, r)
	}

	return reports, nil
}

func ReviewReport(review models.ReportReview) (*models.ReportReview, error) {
	var teamID int
	err := db.Pool.QueryRow(context.Background(), "UPDATE reports SET status = $2 WHERE id = $1 RETURNING team_id", review.ReportID, review.Status).Scan(&teamID)
	if err != nil {
		fmt.Println("Ошибка при обновлении отчета:", err)
		return nil, fmt.Errorf("ошибка при принятии отчета: %w", err)
	}

	if review.Points != nil && *review.Points > 0 {
		_, err = db.Pool.Exec(context.Background(), "UPDATE users SET points = COALESCE(points, 0) + $1 FROM reports WHERE reports.id = $2 AND users.id = reports.author_id", *review.Points, review.ReportID)
		if err != nil {
			fmt.Println("Ошибка при начислении очков пользователю:", err)
			return nil, fmt.Errorf("ошибка при начислении очков пользователю: %w", err)
		}

		_, err = db.Pool.Exec(context.Background(), "UPDATE teams SET points = COALESCE(points, 0) + $1 WHERE id = $2", *review.Points, teamID)
		if err != nil {
			fmt.Println("Ошибка при начислении очков команде:", err)
			return nil, fmt.Errorf("ошибка при начислении очков команде: %w", err)
		}
	}

	var existingReviewID int
	err = db.Pool.QueryRow(context.Background(), "SELECT id FROM report_reviews WHERE report_id = $1 AND reviewer_id = $2", review.ReportID, review.ReviewerID).Scan(&existingReviewID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("ошибка при проверке существующего вердикта: %w", err)
	}

	if existingReviewID > 0 {
		_, err = db.Pool.Exec(context.Background(), "UPDATE report_reviews SET review_text = $1, status = $2 WHERE id = $3", review.ReviewText, review.Status, existingReviewID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обновлении вердикта: %w", err)
		}
	} else {
		row := db.Pool.QueryRow(context.Background(), "INSERT INTO report_reviews (report_id, reviewer_id, review_text, status) VALUES ($1, $2, $3, $4) RETURNING reviewer_id, review_text, status", review.ReportID, review.ReviewerID, review.ReviewText, review.Status)
		if err = row.Scan(&review.ReviewerID, &review.ReviewText, &review.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сохранении вердикта: %w", err)
		}
	}

	return &review, nil
}

// TODO: Refactor this function, fix bug with report data when admin is not team participant
func ReviewDetails(reportID int, userID uuid.UUID) (*models.ReviewDetails, error) {
	var reviewDetails models.ReviewDetails

	err := db.Pool.QueryRow(context.Background(),
		DetailsQuery, reportID, userID).Scan(&reviewDetails.ReviewerID, &reviewDetails.ReviewerUsername,
		&reviewDetails.ReviewText, &reviewDetails.ReportData.ID, &reviewDetails.ReportData.Author,
		&reviewDetails.ReportData.Title, &reviewDetails.ReportData.Description,
		&reviewDetails.ReportData.Status, &reviewDetails.ReportData.Category)

	if errors.Is(err, pgx.ErrNoRows) {
		reportData, reportErr := getReportData(strconv.Itoa(reportID), userID)
		if reportErr != nil {
			fmt.Printf("Error fetching report data: %v", reportErr)
			return nil, reportErr
		}
		return &models.ReviewDetails{
			ReportData: *reportData,
			ReportID:   reportID,
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при получении детального вердикта: %w", err)
	}

	reviewDetails.ReportID = reportID
	return &reviewDetails, nil
}

func getReportData(reportID string, userID uuid.UUID) (*models.ReportData, error) {
	var reportData models.ReportData

	err := db.Pool.QueryRow(context.Background(),
		ReportsQuery, reportID, userID).Scan(&reportData.ID, &reportData.Author, &reportData.Category,
		&reportData.Title, &reportData.Status, &reportData.Description)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных по отчету: %w", err)
	}

	return &reportData, nil
}
