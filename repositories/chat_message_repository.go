package repositories

import (
	"database/sql"
	"fmt"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"
)

func CreateChatMessage(message models.ChatMessage) (models.ChatMessage, error) {
	query := `
		INSERT INTO chat_messages (session_id, role, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := db.DB.QueryRow(
		query,
		message.SessionID,
		message.Role,
		message.Content,
	).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return models.ChatMessage{}, fmt.Errorf("error creating chat message: %v", err)
	}

	return message, nil
}

func GetChatMessagesBySessionID(sessionID string) ([]models.ChatMessage, error) {
	query := `SELECT id, session_id, role, content, created_at 
			  FROM chat_messages WHERE session_id = $1 ORDER BY created_at ASC`

	rows, err := db.DB.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying chat messages: %v", err)
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat message: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func GetChatMessageByID(id string) (models.ChatMessage, error) {
	query := `SELECT id, session_id, role, content, created_at 
			  FROM chat_messages WHERE id = $1`

	var message models.ChatMessage
	err := db.DB.QueryRow(query, id).Scan(
		&message.ID,
		&message.SessionID,
		&message.Role,
		&message.Content,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.ChatMessage{}, fmt.Errorf("chat message not found")
		}
		return models.ChatMessage{}, fmt.Errorf("error querying chat message: %v", err)
	}

	return message, nil
}

func DeleteChatMessage(id string) error {
	query := `DELETE FROM chat_messages WHERE id = $1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting chat message: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat message not found")
	}

	return nil
}
