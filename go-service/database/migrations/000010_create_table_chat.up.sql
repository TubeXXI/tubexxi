CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'text' CHECK (type IN ('text', 'image', 'video', 'audio', 'file', 'location', 'contact')),
    file_url TEXT,
    file_name TEXT,
    file_size BIGINT,
    mime_type VARCHAR(100),
    is_read BOOLEAN DEFAULT FALSE,
    is_delivered BOOLEAN DEFAULT FALSE,
    is_edited BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    reply_to_id UUID REFERENCES chats(id) ON DELETE SET NULL,
    forwarded_from_id UUID REFERENCES chats(id) ON DELETE SET NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure sender and receiver are different users
    CONSTRAINT different_users CHECK (sender_id != receiver_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_chats_sender_id ON chats (sender_id);
CREATE INDEX IF NOT EXISTS idx_chats_receiver_id ON chats (receiver_id);
CREATE INDEX IF NOT EXISTS idx_chats_sender_receiver ON chats (sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_chats_created_at ON chats (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chats_is_read ON chats (is_read);
CREATE INDEX IF NOT EXISTS idx_chats_reply_to_id ON chats (reply_to_id);
CREATE INDEX IF NOT EXISTS idx_chats_forwarded_from_id ON chats (forwarded_from_id);
CREATE INDEX IF NOT EXISTS idx_chats_type ON chats (type);
CREATE INDEX IF NOT EXISTS idx_chats_metadata ON chats USING gin (metadata);

-- Composite index for conversation history
CREATE INDEX IF NOT EXISTS idx_chats_conversation ON chats (sender_id, receiver_id, created_at DESC);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_chats_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_chats_modtime ON chats;
CREATE TRIGGER update_chats_modtime
    BEFORE UPDATE ON chats
    FOR EACH ROW
    EXECUTE FUNCTION update_chats_modtime();

-- Function to get conversation between two users
CREATE OR REPLACE FUNCTION get_conversation(user1_id UUID, user2_id UUID, limit_count INTEGER DEFAULT 50)
RETURNS TABLE (
    id UUID,
    sender_id UUID,
    receiver_id UUID,
    message TEXT,
    type VARCHAR,
    file_url TEXT,
    is_read BOOLEAN,
    is_delivered BOOLEAN,
    is_edited BOOLEAN,
    reply_to_id UUID,
    forwarded_from_id UUID,
    created_at TIMESTAMP WITH TIME ZONE,
    sender_name TEXT,
    sender_avatar TEXT,
    receiver_name TEXT,
    receiver_avatar TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        c.sender_id,
        c.receiver_id,
        c.message,
        c.type,
        c.file_url,
        c.is_read,
        c.is_delivered,
        c.is_edited,
        c.reply_to_id,
        c.forwarded_from_id,
        c.created_at,
        sender.name AS sender_name,
        sender.avatar_url AS sender_avatar,
        receiver.name AS receiver_name,
        receiver.avatar_url AS receiver_avatar
    FROM chats c
    JOIN users sender ON c.sender_id = sender.id
    JOIN users receiver ON c.receiver_id = receiver.id
    WHERE (c.sender_id = user1_id AND c.receiver_id = user2_id)
       OR (c.sender_id = user2_id AND c.receiver_id = user1_id)
    ORDER BY c.created_at DESC
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;