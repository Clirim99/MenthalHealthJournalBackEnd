package models

import (
	"database/sql"
	"fmt"
	"time"
)

type ContextType string

const (
	ContextTypeGlobal     ContextType = "global"
	ContextTypeSingleEntry ContextType = "single_entry"
)

type ChatSession struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	ContextType ContextType `json:"context_type"`
	EntryID     *string     `json:"entry_id,omitempty"` // Only used if context_type is 'single_entry'
	SessionName *string     `json:"session_name,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func CreateChatSessionsTable(db *sql.DB) error {
	// Create enum type
	_, err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE context_type AS ENUM ('global', 'single_entry');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("could not create context_type enum: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS chat_sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		context_type context_type NOT NULL DEFAULT 'global',
		entry_id UUID REFERENCES entries(id) ON DELETE SET NULL,
		session_name VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create chat_sessions table: %v", err)
	}

	_, err = db.Exec(`ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS session_name VARCHAR(100)`)
	if err != nil {
		return fmt.Errorf("could not add session_name to chat_sessions: %v", err)
	}

	// Create indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS chat_sessions_user_id_idx ON chat_sessions(user_id)")
	if err != nil {
		return fmt.Errorf("could not create user_id index: %v", err)
	}

	return nil
}
