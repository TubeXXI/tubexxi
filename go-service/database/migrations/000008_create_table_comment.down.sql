BEGIN;

DROP TRIGGER IF EXISTS update_comments_modtime ON comments;

DROP FUNCTION IF EXISTS update_comments_modtime();

DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_email;
DROP INDEX IF EXISTS idx_comments_page_url;
DROP INDEX IF EXISTS idx_comments_reply_to_id;
DROP INDEX IF EXISTS idx_comments_created_at;

DROP TABLE IF EXISTS comments CASCADE;

COMMIT;