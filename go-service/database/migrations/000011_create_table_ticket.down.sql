BEGIN;

DROP TRIGGER IF EXISTS update_tickets_modtime ON tickets;
DROP FUNCTION IF EXISTS update_tickets_modtime();

DROP INDEX IF EXISTS idx_tickets_user_id;
DROP INDEX IF EXISTS idx_tickets_status;
DROP INDEX IF EXISTS idx_tickets_created_at;

DROP TABLE IF EXISTS tickets CASCADE;

COMMIT;
