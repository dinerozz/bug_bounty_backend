package conversation

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
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

func GetMessages(reportID string, userID uuid.UUID) ([]models.Conversation, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT rc.id, rc.report_id, rc.user_id, rc.message from report_conversations rc JOIN reports r on rc.report_id = r.id JOIN users u ON r.author_id = u.id JOIN team_members tm ON u.id = tm.user_id WHERE report_id = $1", reportID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении сообщений: %w", err)
	}

	var messages []models.Conversation
	for rows.Next() {
		var m models.Conversation
		if err = rows.Scan(&m.ID, &m.ReportID, &m.UserID, &m.Message); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании сообщений: %w", err)
		}
		messages = append(messages, m)
	}

	return messages, nil
}
