CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL,
    scope TEXT NOT NULL DEFAULT 'default',
    value TEXT,
    description TEXT,
    group_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_settings_scope_key ON settings (scope, key);
CREATE INDEX IF NOT EXISTS idx_settings_scope ON settings (scope);

CREATE OR REPLACE FUNCTION update_settings_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_settings_modtime
BEFORE UPDATE ON settings
FOR EACH ROW
EXECUTE FUNCTION update_settings_modtime();
