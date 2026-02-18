package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pgvector/pgvector-go"
)

type Entry struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	Content       string         `json:"content"`
	Embedding     pgvector.Vector `json:"-"` // Not exposed in JSON, but stored in DB
	SentimentScore int           `json:"sentiment_score"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func CreateEntriesTable(db *sql.DB) error {
	// Enable pgvector extension
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return fmt.Errorf("could not enable vector extension: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS entries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		embedding vector(1536),
		sentiment_score INTEGER CHECK (sentiment_score >= 1 AND sentiment_score <= 10),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create entries table: %v", err)
	}

	// Create vector index for similarity search
	indexQuery := `
	CREATE INDEX IF NOT EXISTS entries_embedding_idx ON entries 
	USING ivfflat (embedding vector_cosine_ops)
	WITH (lists = 100);`

	_, err = db.Exec(indexQuery)
	if err != nil {
		// Index creation might fail if table is empty, that's okay
		fmt.Printf("Warning: Could not create vector index (this is normal if table is empty): %v\n", err)
	}

	// Create other indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS entries_user_id_idx ON entries(user_id)")
	if err != nil {
		return fmt.Errorf("could not create user_id index: %v", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS entries_created_at_idx ON entries(created_at)")
	if err != nil {
		return fmt.Errorf("could not create created_at index: %v", err)
	}

	return nil
}
