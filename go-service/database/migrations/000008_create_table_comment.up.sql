CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    page_url TEXT NOT NULL,
    name TEXT,
    email TEXT,
    comment TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'text' CHECK (type IN ('text', 'image', 'video', 'link', 'document', 'other')),
    reply_to_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    is_edited BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    media_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_email ON comments (email);
CREATE INDEX IF NOT EXISTS idx_comments_page_url ON comments (page_url);
CREATE INDEX IF NOT EXISTS idx_comments_reply_to_id ON comments (reply_to_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_comments_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_comments_modtime ON comments;
CREATE TRIGGER update_comments_modtime
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_comments_modtime();