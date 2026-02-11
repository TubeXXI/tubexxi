CREATE TABLE IF NOT EXISTS chat_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(20) DEFAULT 'like' CHECK (reaction_type IN ('like', 'love', 'haha', 'wow', 'sad', 'angry', 'thumbs_up', 'thumbs_down')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(chat_id, user_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_chat_likes_chat_id ON chat_likes (chat_id);
CREATE INDEX IF NOT EXISTS idx_chat_likes_user_id ON chat_likes (user_id);
CREATE INDEX IF NOT EXISTS idx_chat_likes_reaction_type ON chat_likes (reaction_type);
CREATE INDEX IF NOT EXISTS idx_chat_likes_created_at ON chat_likes (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_likes_chat_reaction ON chat_likes (chat_id, reaction_type);

-- Function to get reactions for a chat
CREATE OR REPLACE FUNCTION get_chat_reactions(chat_uuid UUID)
RETURNS TABLE (
    reaction_type VARCHAR,
    count BIGINT,
    users JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        cl.reaction_type,
        COUNT(*) AS count,
        jsonb_agg(jsonb_build_object(
            'user_id', cl.user_id,
            'name', u.name,
            'avatar', u.avatar_url,
            'created_at', cl.created_at
        )) AS users
    FROM chat_likes cl
    JOIN users u ON cl.user_id = u.id
    WHERE cl.chat_id = chat_uuid
    GROUP BY cl.reaction_type
    ORDER BY count DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get user's reaction to a specific chat
CREATE OR REPLACE FUNCTION get_user_chat_reaction(chat_uuid UUID, user_uuid UUID)
RETURNS TABLE (
    reaction_type VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT cl.reaction_type, cl.created_at
    FROM chat_likes cl
    WHERE cl.chat_id = chat_uuid AND cl.user_id = user_uuid
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;