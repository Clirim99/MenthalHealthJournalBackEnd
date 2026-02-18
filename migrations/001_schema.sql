-- Enable pgvector extension for vector similarity search
CREATE EXTENSION IF NOT EXISTS vector;

-- Users table (already exists, but included for reference)
-- Assuming users table uses UUID as primary key

-- Entries table with vector embeddings
CREATE TABLE IF NOT EXISTS entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding vector(1536), -- OpenAI text-embedding-3-small produces 1536-dimensional vectors
    sentiment_score INTEGER CHECK (sentiment_score >= 1 AND sentiment_score <= 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for vector similarity search (cosine distance)
CREATE INDEX IF NOT EXISTS entries_embedding_idx ON entries 
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Create index for user_id lookups
CREATE INDEX IF NOT EXISTS entries_user_id_idx ON entries(user_id);

-- Create index for date-based queries
CREATE INDEX IF NOT EXISTS entries_created_at_idx ON entries(created_at);

-- Summaries table
CREATE TYPE summary_type AS ENUM ('weekly', 'monthly');

CREATE TABLE IF NOT EXISTS summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    type summary_type NOT NULL,
    date_range_start DATE NOT NULL,
    date_range_end DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for summaries
CREATE INDEX IF NOT EXISTS summaries_user_id_idx ON summaries(user_id);
CREATE INDEX IF NOT EXISTS summaries_type_idx ON summaries(type);
CREATE INDEX IF NOT EXISTS summaries_date_range_idx ON summaries(date_range_start, date_range_end);

-- ChatSessions table
CREATE TYPE context_type AS ENUM ('global', 'single_entry');

CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    context_type context_type NOT NULL DEFAULT 'global',
    entry_id UUID REFERENCES entries(id) ON DELETE SET NULL, -- Only used if context_type is 'single_entry'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for chat sessions
CREATE INDEX IF NOT EXISTS chat_sessions_user_id_idx ON chat_sessions(user_id);
CREATE INDEX IF NOT EXISTS chat_sessions_context_type_idx ON chat_sessions(context_type);

-- ChatMessages table
CREATE TYPE message_role AS ENUM ('user', 'assistant');

CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role message_role NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for chat messages
CREATE INDEX IF NOT EXISTS chat_messages_session_id_idx ON chat_messages(session_id);
CREATE INDEX IF NOT EXISTS chat_messages_created_at_idx ON chat_messages(created_at);
