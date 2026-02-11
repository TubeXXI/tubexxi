BEGIN;

DROP TRIGGER IF EXISTS update_chat_rooms_modtime ON chat_rooms;
DROP TRIGGER IF EXISTS update_chat_messages_modtime ON chat_messages;
DROP TRIGGER IF EXISTS update_room_last_message ON chat_messages;

DROP FUNCTION IF EXISTS update_chat_messages_modtime();
DROP FUNCTION IF EXISTS update_room_last_message();
DROP FUNCTION IF EXISTS create_personal_chat_room();

DROP INDEX IF EXISTS idx_chat_rooms_type;
DROP INDEX IF EXISTS idx_chat_rooms_last_message_at;
DROP INDEX IF EXISTS idx_chat_rooms_created_by;
DROP INDEX IF EXISTS idx_room_participants_room_id;
DROP INDEX IF EXISTS idx_room_participants_user_id;
DROP INDEX IF EXISTS idx_room_participants_last_read;
DROP INDEX IF EXISTS idx_chat_messages_room_id;
DROP INDEX IF EXISTS idx_chat_messages_sender_id;
DROP INDEX IF EXISTS idx_chat_messages_created_at;
DROP INDEX IF EXISTS idx_chat_messages_reply_to;
DROP INDEX IF EXISTS idx_chat_messages_type;
DROP INDEX IF EXISTS idx_chat_messages_is_pinned;
DROP INDEX IF EXISTS idx_chat_messages_metadata;
DROP INDEX IF EXISTS idx_message_reactions_message_id;
DROP INDEX IF EXISTS idx_message_reactions_user_id;

DROP TABLE IF EXISTS chat_rooms CASCADE;
DROP TABLE IF EXISTS chat_room_participants CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_message_reactions CASCADE;
DROP TABLE IF EXISTS chat_message_reads CASCADE;
DROP TABLE IF EXISTS chat_message_deliveries CASCADE;

COMMIT;