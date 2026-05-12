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
		return pgvector.Vector{}, fmt.Errorf("error creating embedding: %v", err) // Changed nil to pgvector.Vector{}
	}

	if len(resp.Data) == 0 {
		return pgvector.Vector{}, fmt.Errorf("no embedding data returned") // Changed nil to pgvector.Vector{}
	}

	// pgvector.NewVector expects []float32, and OpenAI already returns []float32!
	// No conversion loop is needed. Just pass it directly.
	return pgvector.NewVector(resp.Data[0].Embedding), nil
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

// GetChatSessionTitleLLM generates a short session title using the same model as chat (GPT-4o).
func GetChatSessionTitleLLM(firstUserMessage string) (string, error) {
	if OpenAIClient == nil {
		InitOpenAIClient()
	}

	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You write short chat thread titles for a mental health journaling app. Given the user's first message, reply with ONLY a concise title (maximum 100 characters), no quotation marks, no prefix or explanation. Capture the main topic or emotional theme.`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: firstUserMessage,
			},
		},
		Temperature: 0.4,
		MaxTokens:   80,
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
