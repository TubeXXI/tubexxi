CREATE TABLE IF NOT EXISTS view_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT,
    page_url TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    browser_language TEXT,
    device_type TEXT,
    platform TEXT NOT NULL DEFAULT 'web' CHECK (platform IN ('web', 'mobile')),
    view_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(50) NOT NULL DEFAULT 'movies' CHECK (type IN ('movies', 'series', 'animes')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_view_history_user_id ON view_histories (user_id);
CREATE INDEX IF NOT EXISTS idx_view_history_page_url ON view_histories (page_url);
CREATE INDEX IF NOT EXISTS idx_view_history_type ON view_histories (type);



CREATE OR REPLACE FUNCTION update_view_histories_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_view_histories_modtime
BEFORE UPDATE ON view_histories
FOR EACH ROW
EXECUTE FUNCTION update_view_histories_modtime();
