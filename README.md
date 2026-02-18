# Mental Health AI Journal Backend

A Go-based backend for a mental health journaling application with AI-powered analysis using RAG (Retrieval Augmented Generation).

## Features

- **Journal Entries**: Store journal entries with automatic embedding generation
- **RAG Chat**: Chat with your journal using AI that retrieves relevant past entries
- **Vector Search**: PostgreSQL with pgvector for semantic similarity search
- **Chat Sessions**: Maintain conversation history with your journal

## Tech Stack

- **Backend**: Go (Golang) with Gin framework
- **Database**: PostgreSQL with pgvector extension
- **AI**: OpenAI API (GPT-4o for chat, text-embedding-3-small for embeddings)

## Prerequisites

1. **PostgreSQL** with pgvector extension installed
   ```bash
   # Install pgvector extension
   # For macOS: brew install pgvector
   # Or follow: https://github.com/pgvector/pgvector#installation
   ```

2. **Go** 1.24.2 or later

3. **OpenAI API Key**

## Setup

### 1. Database Setup

Create a PostgreSQL database and enable the pgvector extension:

```sql
CREATE DATABASE MenthalHealthCare;
\c MenthalHealthCare;
CREATE EXTENSION IF NOT EXISTS vector;
```

Or run the migration file:
```bash
psql -U postgres -d MenthalHealthCare -f migrations/001_schema.sql
```

### 2. Environment Variables

Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key-here"
```

### 3. Database Configuration

Update the database connection in `db/db.go`:
```go
host     = "localhost"
port     = 5432
user     = "postgres"
password = "root"
dbname   = "MenthalHealthCare"
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Entries

- **POST /entries** - Create a new journal entry
  ```json
  {
    "user_id": "uuid",
    "content": "Today I felt...",
    "sentiment_score": 7
  }
  ```

- **GET /entries/:id** - Get a specific entry
- **GET /users/:user_id/entries** - Get all entries for a user
- **PUT /entries/:id** - Update an entry
- **DELETE /entries/:id** - Delete an entry

### Chat

- **POST /chat** - Chat with your journal (RAG pipeline)
  ```json
  {
    "user_id": "uuid",
    "message": "Why was I sad last week?",
    "session_id": "optional-session-id"
  }
  ```

- **POST /chat/sessions** - Create a new chat session
- **GET /chat/sessions/:id** - Get a chat session
- **GET /users/:user_id/chat/sessions** - Get all sessions for a user
- **GET /chat/sessions/:session_id/messages** - Get chat history
- **DELETE /chat/sessions/:id** - Delete a chat session

## How RAG Works

1. **User sends a message** (e.g., "Why was I sad last week?")
2. **Vectorize**: The message is converted to an embedding using OpenAI
3. **Search**: PostgreSQL performs cosine similarity search to find the top 5 most relevant journal entries
4. **Context Building**: Relevant entries are formatted into a context string
5. **Generate**: GPT-4o generates a response using the context and user's question
6. **Return**: The AI response is returned to the client

## Database Schema

### Entries
- `id` (UUID)
- `user_id` (UUID, FK to users)
- `content` (TEXT)
- `embedding` (vector(1536))
- `sentiment_score` (INT, 1-10)
- `created_at`, `updated_at` (TIMESTAMP)

### Chat Sessions
- `id` (UUID)
- `user_id` (UUID, FK to users)
- `context_type` (ENUM: 'global', 'single_entry')
- `entry_id` (UUID, optional, FK to entries)
- `created_at`, `updated_at` (TIMESTAMP)

### Chat Messages
- `id` (UUID)
- `session_id` (UUID, FK to chat_sessions)
- `role` (ENUM: 'user', 'assistant')
- `content` (TEXT)
- `created_at` (TIMESTAMP)

## Project Structure

```
.
├── controllers/     # HTTP handlers
├── db/             # Database connection
├── migrations/      # SQL migration files
├── models/         # Data models
├── repositories/   # Database operations
├── router/         # Route definitions
├── services/       # Business logic (including RAG)
└── main.go        # Application entry point
```

## Key Functions

### `GetAnswerFromJournal(userID, userMessage string)`
The core RAG pipeline function that:
1. Converts user message to embedding
2. Searches for similar entries
3. Builds context
4. Generates AI response

Located in: `services/rag_service.go`

## Notes

- Embeddings are automatically generated when creating entries
- The vector index (ivfflat) is created automatically but may need data to be effective
- Chat sessions maintain conversation history
- All embeddings are stored in PostgreSQL using pgvector

## Troubleshooting

1. **pgvector not found**: Make sure the extension is installed in PostgreSQL
2. **OpenAI errors**: Check your API key and rate limits
3. **Vector index warnings**: Normal if the table is empty, will work once entries are added
