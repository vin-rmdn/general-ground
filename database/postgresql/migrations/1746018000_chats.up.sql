CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chats_created_at ON messages (created_at);

-- from_user_id index is intentionally not created because it is covered by the composite index below
CREATE INDEX IF NOT EXISTS idx_chats_to_user_id ON messages (to_user_id);
CREATE INDEX IF NOT EXISTS idx_chats_from_to_user_id ON messages (from_user_id, to_user_id);
