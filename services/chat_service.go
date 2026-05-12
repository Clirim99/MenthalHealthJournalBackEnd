package services

import (
	"log"
	"net/http"
	"menthalhealthjournal/models"
	"menthalhealthjournal/repositories"

	"github.com/gin-gonic/gin"
)

func CreateChatSession(c *gin.Context) {
	var req struct {
		UserID      string  `json:"user_id" binding:"required"`
		ContextType string  `json:"context_type"` // "global" or "single_entry"
		EntryID     *string `json:"entry_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	session := models.ChatSession{
		UserID:      req.UserID,
		ContextType: models.ContextTypeGlobal,
		EntryID:     req.EntryID,
	}

	if req.ContextType == "single_entry" {
		session.ContextType = models.ContextTypeSingleEntry
		if req.EntryID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "entry_id is required for single_entry context type"})
			return
		}
	}

	createdSession, err := repositories.CreateChatSession(session)
	if err != nil {
		log.Println("Error creating chat session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdSession)
}

func GetChatSession(c *gin.Context) {
	id := c.Param("id")

	session, err := repositories.GetChatSessionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func GetChatSessionsByUser(c *gin.Context) {
	userID := c.Param("user_id")

	sessions, err := repositories.GetChatSessionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func ChatWithJournal(c *gin.Context) {
	var req struct {
		Message   string `json:"message" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
		SessionID string `json:"session_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	var response string
	var err error

	if req.SessionID != "" {
		// Use existing session
		response, err = GetAnswerFromJournalWithSession(req.SessionID, req.Message)
		if err != nil {
			log.Println("Error in RAG pipeline:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response":    response,
			"session_id": req.SessionID,
		})
	} else {
		// Create new session and use RAG pipeline
		response, err = GetAnswerFromJournal(req.UserID, req.Message)
		if err != nil {
			log.Println("Error in RAG pipeline:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create a new session for this conversation
		session := models.ChatSession{
			UserID:      req.UserID,
			ContextType: models.ContextTypeGlobal,
		}
		createdSession, err := repositories.CreateChatSession(session)
		sessionID := ""
		if err != nil {
			log.Println("Error creating chat session:", err)
			// Continue anyway, just return the response without session
		} else {
			sessionID = createdSession.ID
			// Save messages to session
			userMsg := models.ChatMessage{
				SessionID: createdSession.ID,
				Role:      models.MessageRoleUser,
				Content:   req.Message,
			}
			if _, err := repositories.CreateChatMessage(userMsg); err != nil {
				log.Println("Error saving user message:", err)
			} else {
				MaybeSetSessionNameOnFirstUserMessage(createdSession, req.Message)
			}

			assistantMsg := models.ChatMessage{
				SessionID: createdSession.ID,
				Role:      models.MessageRoleAssistant,
				Content:   response,
			}
			repositories.CreateChatMessage(assistantMsg)
		}

		c.JSON(http.StatusOK, gin.H{
			"response":    response,
			"session_id": sessionID,
		})
	}
}

func GetChatHistory(c *gin.Context) {
	sessionID := c.Param("id")
	
	messages, err := repositories.GetChatMessagesBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func DeleteChatSession(c *gin.Context) {
	id := c.Param("id")

	err := repositories.DeleteChatSession(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat session deleted successfully"})
}
