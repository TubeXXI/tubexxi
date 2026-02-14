CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL,
    package_name TEXT NOT NULL DEFAULT 'default',
    value TEXT,
    description TEXT,
    group_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_applications_package_name_key ON applications (package_name, key);
CREATE INDEX IF NOT EXISTS idx_applications_package_name ON applications (package_name);

CREATE OR REPLACE FUNCTION update_applications_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_applications_modtime
BEFORE UPDATE ON applications
FOR EACH ROW
EXECUTE FUNCTION update_applications_modtime();
