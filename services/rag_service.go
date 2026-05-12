package services

import (
	"fmt"
	"menthalhealthjournal/models"
	"menthalhealthjournal/repositories"
	//"time"
)

// GetAnswerFromJournal implements the RAG pipeline for chatting with the journal
// This is the core function that handles the entire RAG workflow
func GetAnswerFromJournal(userID, userMessage string) (string, error) {
	// Step 1: Vectorize the user's message
	queryEmbedding, err := GetEmbedding(userMessage)
	if err != nil {
		return "", fmt.Errorf("error creating embedding for user message: %v", err)
	}

	// Step 2: Perform cosine similarity search to find top 5 most similar entries
	similarEntries, err := repositories.FindSimilarEntries(userID, queryEmbedding, 5)
	if err != nil {
		return "", fmt.Errorf("error finding similar entries: %v", err)
	}

	// Step 3: Format the found entries into a context string
	contextString := buildContextFromEntries(similarEntries)

	// Step 4: Construct the system prompt
	systemPrompt := `You are an empathetic mental health assistant. Answer the user's question using strictly the following context from their journal entries. If the answer is not in the context, state that you don't have enough information from their journal to answer that question. Be supportive, understanding, and focus on helping them understand patterns in their mental health journey.`

	// Step 5: Generate response using GPT-4o
	response, err := GetChatCompletion(systemPrompt, userMessage+"\n\nContext from journal:\n"+contextString)
	if err != nil {
		return "", fmt.Errorf("error generating chat completion: %v", err)
	}

	return response, nil
}

// buildContextFromEntries formats journal entries into a readable context string
func buildContextFromEntries(entries []models.Entry) string {
	if len(entries) == 0 {
		return "No relevant journal entries found."
	}

	context := ""
	for _, entry := range entries {
		dateStr := entry.CreatedAt.Format("2006-01-02")
		context += fmt.Sprintf("- [%s]: %s\n", dateStr, entry.Content)
	}

	return context
}

// GetAnswerFromJournalWithSession handles RAG with chat session management
func GetAnswerFromJournalWithSession(sessionID, userMessage string) (string, error) {
	// Get the chat session to determine context type
	session, err := repositories.GetChatSessionByID(sessionID)
	if err != nil {
		return "", fmt.Errorf("error getting chat session: %v", err)
	}

	// Save user message to chat history
	userMsg := models.ChatMessage{
		SessionID: sessionID,
		Role:      models.MessageRoleUser,
		Content:   userMessage,
	}
	_, err = repositories.CreateChatMessage(userMsg)
	if err != nil {
		return "", fmt.Errorf("error saving user message: %v", err)
	}

	MaybeSetSessionNameOnFirstUserMessage(session, userMessage)

	var response string

	// Handle different context types
	if session.ContextType == models.ContextTypeSingleEntry && session.EntryID != nil {
		// Single entry context - get that specific entry
		entry, err := repositories.GetEntryByID(*session.EntryID)
		if err != nil {
			return "", fmt.Errorf("error getting entry: %v", err)
		}

		contextString := fmt.Sprintf("- [%s]: %s\n", entry.CreatedAt.Format("2006-01-02"), entry.Content)
		systemPrompt := `You are an empathetic mental health assistant. Answer the user's question using strictly the following context from their journal entry. Be supportive and understanding.`

		response, err = GetChatCompletion(systemPrompt, userMessage+"\n\nContext from journal:\n"+contextString)
		if err != nil {
			return "", fmt.Errorf("error generating chat completion: %v", err)
		}
	} else {
		// Global context - use RAG pipeline
		response, err = GetAnswerFromJournal(session.UserID, userMessage)
		if err != nil {
			return "", err
		}
	}

	// Save assistant response to chat history
	assistantMsg := models.ChatMessage{
		SessionID: sessionID,
		Role:      models.MessageRoleAssistant,
		Content:   response,
	}
	_, err = repositories.CreateChatMessage(assistantMsg)
	if err != nil {
		return "", fmt.Errorf("error saving assistant message: %v", err)
	}

	return response, nil
}
