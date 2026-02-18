package services

import (
	"context"
	"fmt"
	"os"

	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func InitOpenAIClient() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}
	OpenAIClient = openai.NewClient(apiKey)
}

// GetEmbedding converts text to a vector embedding using OpenAI's text-embedding-3-small model
func GetEmbedding(text string) (pgvector.Vector, error) {
	if OpenAIClient == nil {
		InitOpenAIClient()
	}

	ctx := context.Background()
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	}

	resp, err := OpenAIClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating embedding: %v", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// Convert []float32 to []float64 for pgvector
	embedding := make([]float64, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float64(v)
	}

	return pgvector.NewVector(embedding), nil
}

// GetChatCompletion sends a message to GPT-4o and returns the response
func GetChatCompletion(systemPrompt, userMessage string) (string, error) {
	if OpenAIClient == nil {
		InitOpenAIClient()
	}

	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage,
			},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	resp, err := OpenAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
