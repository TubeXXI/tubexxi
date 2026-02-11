BEGIN;

DROP TRIGGER IF EXISTS update_view_histories_modtime ON view_histories;

DROP FUNCTION IF EXISTS update_view_histories_modtime();

DROP INDEX IF EXISTS idx_view_history_user_id;
DROP INDEX IF EXISTS idx_view_history_page_url;

DROP TABLE IF EXISTS view_histories;

COMMIT;