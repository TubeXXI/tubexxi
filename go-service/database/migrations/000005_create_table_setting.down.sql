BEGIN;

DROP TRIGGER IF EXISTS update_settings_modtime ON settings;

DROP FUNCTION IF EXISTS update_settings_modtime();

DROP INDEX IF EXISTS idx_settings_scope;
DROP INDEX IF EXISTS idx_settings_scope_key;

DROP TABLE IF EXISTS settings;

COMMIT;