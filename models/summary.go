package models

import (
	"database/sql"
	"fmt"
	"time"
)

type SummaryType string

const (
	SummaryTypeWeekly  SummaryType = "weekly"
	SummaryTypeMonthly SummaryType = "monthly"
)

type Summary struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Content       string    `json:"content"`
	Type          SummaryType `json:"type"`
	DateRangeStart time.Time `json:"date_range_start"`
	DateRangeEnd   time.Time `json:"date_range_end"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func CreateSummariesTable(db *sql.DB) error {
	// Create enum type
	_, err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE summary_type AS ENUM ('weekly', 'monthly');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("could not create summary_type enum: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS summaries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		type summary_type NOT NULL,
		date_range_start DATE NOT NULL,
		date_range_end DATE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create summaries table: %v", err)
	}

	// Create indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS summaries_user_id_idx ON summaries(user_id)")
	if err != nil {
		return fmt.Errorf("could not create user_id index: %v", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS summaries_type_idx ON summaries(type)")
	if err != nil {
		return fmt.Errorf("could not create type index: %v", err)
	}

	return nil
}
