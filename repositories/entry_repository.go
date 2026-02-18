package repositories

import (
	"database/sql"
	"fmt"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"

	"github.com/pgvector/pgvector-go"
)

func CreateEntry(entry models.Entry) (models.Entry, error) {
	query := `
		INSERT INTO entries (user_id, content, embedding, sentiment_score)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := db.DB.QueryRow(
		query,
		entry.UserID,
		entry.Content,
		entry.Embedding,
		entry.SentimentScore,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)

	if err != nil {
		return models.Entry{}, fmt.Errorf("error creating entry: %v", err)
	}

	return entry, nil
}

func GetEntryByID(id string) (models.Entry, error) {
	query := `SELECT id, user_id, content, embedding, sentiment_score, created_at, updated_at 
			  FROM entries WHERE id = $1`

	var entry models.Entry
	var embedding pgvector.Vector
	err := db.DB.QueryRow(query, id).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.Content,
		&embedding,
		&entry.SentimentScore,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}, fmt.Errorf("entry not found")
		}
		return models.Entry{}, fmt.Errorf("error querying entry: %v", err)
	}

	entry.Embedding = embedding
	return entry, nil
}

func GetEntriesByUserID(userID string, limit int) ([]models.Entry, error) {
	query := `SELECT id, user_id, content, embedding, sentiment_score, created_at, updated_at 
			  FROM entries WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := db.DB.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying entries: %v", err)
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		var embedding pgvector.Vector
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Content,
			&embedding,
			&entry.SentimentScore,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning entry: %v", err)
		}
		entry.Embedding = embedding
		entries = append(entries, entry)
	}

	return entries, nil
}

// FindSimilarEntries performs cosine similarity search to find the top N most similar entries
func FindSimilarEntries(userID string, queryEmbedding pgvector.Vector, limit int) ([]models.Entry, error) {
	// Using cosine distance (1 - cosine similarity)
	// Lower distance = higher similarity
	query := `
		SELECT id, user_id, content, embedding, sentiment_score, created_at, updated_at,
		       1 - (embedding <=> $1) as similarity
		FROM entries 
		WHERE user_id = $2 AND embedding IS NOT NULL
		ORDER BY embedding <=> $1
		LIMIT $3
	`

	rows, err := db.DB.Query(query, queryEmbedding, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("error performing similarity search: %v", err)
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		var embedding pgvector.Vector
		var similarity float64
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Content,
			&embedding,
			&entry.SentimentScore,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning entry: %v", err)
		}
		entry.Embedding = embedding
		entries = append(entries, entry)
	}

	return entries, nil
}

func UpdateEntry(id string, entry models.Entry) (models.Entry, error) {
	query := `
		UPDATE entries 
		SET content = $1, embedding = $2, sentiment_score = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, user_id, content, embedding, sentiment_score, created_at, updated_at
	`

	var updatedEntry models.Entry
	var embedding pgvector.Vector
	err := db.DB.QueryRow(
		query,
		entry.Content,
		entry.Embedding,
		entry.SentimentScore,
		id,
	).Scan(
		&updatedEntry.ID,
		&updatedEntry.UserID,
		&updatedEntry.Content,
		&embedding,
		&updatedEntry.SentimentScore,
		&updatedEntry.CreatedAt,
		&updatedEntry.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}, fmt.Errorf("entry not found")
		}
		return models.Entry{}, fmt.Errorf("error updating entry: %v", err)
	}

	updatedEntry.Embedding = embedding
	return updatedEntry, nil
}

func DeleteEntry(id string) error {
	query := `DELETE FROM entries WHERE id = $1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting entry: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("entry not found")
	}

	return nil
}
