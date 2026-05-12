package repositories

import (
	"database/sql"
	"fmt"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"
)

func CreateChatSession(session models.ChatSession) (models.ChatSession, error) {
	query := `
		INSERT INTO chat_sessions (user_id, context_type, entry_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var entryID interface{}
	if session.EntryID != nil {
		entryID = *session.EntryID
	} else {
		entryID = nil
	}

	err := db.DB.QueryRow(
		query,
		session.UserID,
		session.ContextType,
		entryID,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return models.ChatSession{}, fmt.Errorf("error creating chat session: %v", err)
	}

	return session, nil
}

func GetChatSessionByID(id string) (models.ChatSession, error) {
	query := `SELECT id, user_id, context_type, entry_id, session_name, created_at, updated_at 
			  FROM chat_sessions WHERE id = $1`

	var session models.ChatSession
	var entryID sql.NullString
	var sessionName sql.NullString
	err := db.DB.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.ContextType,
		&entryID,
		&sessionName,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.ChatSession{}, fmt.Errorf("chat session not found")
		}
		return models.ChatSession{}, fmt.Errorf("error querying chat session: %v", err)
	}

	if entryID.Valid {
		session.EntryID = &entryID.String
	}
	if sessionName.Valid {
		s := sessionName.String
		session.SessionName = &s
	}

	return session, nil
}

func GetChatSessionsByUserID(userID string) ([]models.ChatSession, error) {
	query := `SELECT id, user_id, context_type, entry_id, session_name, created_at, updated_at 
			  FROM chat_sessions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chat sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.ChatSession
	for rows.Next() {
		var session models.ChatSession
		var entryID sql.NullString
		var sessionName sql.NullString
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.ContextType,
			&entryID,
			&sessionName,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat session: %v", err)
		}
		if entryID.Valid {
			session.EntryID = &entryID.String
		}
		if sessionName.Valid {
			s := sessionName.String
			session.SessionName = &s
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func UpdateChatSession(id string, session models.ChatSession) error {
	query := `
		UPDATE chat_sessions 
		SET context_type = $1, entry_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	var entryID interface{}
	if session.EntryID != nil {
		entryID = *session.EntryID
	} else {
		entryID = nil
	}

	result, err := db.DB.Exec(query, session.ContextType, entryID, id)
	if err != nil {
		return fmt.Errorf("error updating chat session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	return nil
}

// UpdateChatSessionNameIfUnset sets session_name only when it is currently NULL (first automatic title).
func UpdateChatSessionNameIfUnset(sessionID, name string) error {
	_, err := db.DB.Exec(`
		UPDATE chat_sessions
		SET session_name = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND session_name IS NULL`,
		sessionID, name)
	if err != nil {
		return fmt.Errorf("error updating chat session name: %v", err)
	}
	return nil
}

func DeleteChatSession(id string) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting chat session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	return nil
}
