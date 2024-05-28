package conversation

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/auth"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
)

func SendMessage(conversation models.Conversation) (*models.Conversation, error) {
	// TODO: add anti-spam mechanism
	err := db.Pool.QueryRow(context.Background(), "INSERT INTO report_conversations (report_id, user_id, message) VALUES ($1, $2, $3) returning id, report_id, user_id, message",
		conversation.ReportID, conversation.UserID, conversation.Message).Scan(&conversation.ID, &conversation.ReportID, &conversation.UserID, &conversation.Message)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении сообщения: %w", err)
	}

	return &conversation, nil
}

func GetMessages(reportID string, userID uuid.UUID) ([]models.GetConversation, error) {
	query := `
		SELECT rc.id, rc.report_id, rc.user_id, rc.message, u.username 
		FROM report_conversations rc
		JOIN reports r ON rc.report_id = r.id
		JOIN users u ON rc.user_id = u.id
		JOIN team_members tm ON u.id = tm.user_id
		WHERE rc.report_id = $1
		ORDER BY rc.created_at
	`

	rows, err := db.Pool.Query(context.Background(), query, reportID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении сообщений: %w", err)
	}
	defer rows.Close()

	var messages []models.GetConversation
	for rows.Next() {
		var m models.GetConversation
		if err = rows.Scan(&m.ID, &m.ReportID, &m.UserID, &m.Message, &m.Username); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании сообщений: %w", err)
		}
		isAdmin, validateRoleErr := auth.ValidateRole(m.UserID)
		if validateRoleErr != nil {
			return nil, fmt.Errorf("ошибка при валидации роли пользователя: %w", validateRoleErr)
		}
		m.IsAdmin = isAdmin
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return messages, nil
}

func IsUserInTeam(reportID string, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM team_members tm
		JOIN reports r ON tm.team_id = r.team_id
		WHERE r.id = $1 AND tm.user_id = $2
	`
	var count int
	err := db.Pool.QueryRow(context.Background(), query, reportID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке принадлежности к команде: %w", err)
	}
	return count > 0, nil
}
