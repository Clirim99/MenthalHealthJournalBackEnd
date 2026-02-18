package models

import (
	"database/sql"
	"fmt"
	"time"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type ChatMessage struct {
	ID        string      `json:"id"`
	SessionID string      `json:"session_id"`
	Role      MessageRole `json:"role"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
}

func CreateChatMessagesTable(db *sql.DB) error {
	// Create enum type
	_, err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE message_role AS ENUM ('user', 'assistant');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("could not create message_role enum: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS chat_messages (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
		role message_role NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create chat_messages table: %v", err)
	}

	// Create indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS chat_messages_session_id_idx ON chat_messages(session_id)")
	if err != nil {
		return fmt.Errorf("could not create session_id index: %v", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS chat_messages_created_at_idx ON chat_messages(created_at)")
	if err != nil {
		return fmt.Errorf("could not create created_at index: %v", err)
	}

	return nil
}
