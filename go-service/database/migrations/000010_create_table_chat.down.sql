BEGIN;

DROP FUNCTION IF EXISTS get_conversation(UUID, UUID, INTEGER);
DROP TRIGGER IF EXISTS update_chats_modtime ON chats;
DROP FUNCTION IF EXISTS update_chats_modtime();

DROP INDEX IF EXISTS idx_chats_sender_id;
DROP INDEX IF EXISTS idx_chats_receiver_id;
DROP INDEX IF EXISTS idx_chats_sender_receiver;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP INDEX IF EXISTS idx_chats_is_read;
DROP INDEX IF EXISTS idx_chats_reply_to_id;
DROP INDEX IF EXISTS idx_chats_forwarded_from_id;
DROP INDEX IF EXISTS idx_chats_type;
DROP INDEX IF EXISTS idx_chats_metadata;
DROP INDEX IF EXISTS idx_chats_conversation;

DROP TABLE IF EXISTS chats CASCADE;

COMMIT;