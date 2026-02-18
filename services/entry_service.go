package services

import (
	"log"
	"net/http"
	"menthalhealthjournal/models"
	"menthalhealthjournal/repositories"

	"github.com/gin-gonic/gin"
)

func CreateEntry(c *gin.Context) {
	var req struct {
		Content       string `json:"content" binding:"required"`
		SentimentScore int   `json:"sentiment_score" binding:"required"`
		UserID        string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error decoding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Step 1: Get embedding from OpenAI
	embedding, err := GetEmbedding(req.Content)
	if err != nil {
		log.Println("Error creating embedding:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create embedding: " + err.Error()})
		return
	}

	// Step 2: Create entry with embedding
	entry := models.Entry{
		UserID:        req.UserID,
		Content:       req.Content,
		Embedding:     embedding,
		SentimentScore: req.SentimentScore,
	}

	createdEntry, err := repositories.CreateEntry(entry)
	if err != nil {
		log.Println("Error creating entry:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't expose embedding in response
	createdEntry.Embedding = nil
	c.JSON(http.StatusCreated, createdEntry)
}

func GetEntry(c *gin.Context) {
	id := c.Param("id")

	entry, err := repositories.GetEntryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Don't expose embedding in response
	entry.Embedding = nil
	c.JSON(http.StatusOK, entry)
}

func GetEntriesByUser(c *gin.Context) {
	userID := c.Param("user_id")
	limit := 50 // Default limit

	entries, err := repositories.GetEntriesByUserID(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove embeddings from response
	for i := range entries {
		entries[i].Embedding = nil
	}

	c.JSON(http.StatusOK, entries)
}

func UpdateEntry(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Content       string `json:"content"`
		SentimentScore int   `json:"sentiment_score"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Get existing entry
	existingEntry, err := repositories.GetEntryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Content != "" {
		existingEntry.Content = req.Content
		// Re-generate embedding if content changed
		embedding, err := GetEmbedding(req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create embedding: " + err.Error()})
			return
		}
		existingEntry.Embedding = embedding
	}
	if req.SentimentScore > 0 {
		existingEntry.SentimentScore = req.SentimentScore
	}

	updatedEntry, err := repositories.UpdateEntry(id, existingEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't expose embedding in response
	updatedEntry.Embedding = nil
	c.JSON(http.StatusOK, updatedEntry)
}

func DeleteEntry(c *gin.Context) {
	id := c.Param("id")

	err := repositories.DeleteEntry(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}
