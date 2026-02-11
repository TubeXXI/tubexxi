BEGIN;

DROP FUNCTION IF EXISTS get_chat_reactions(UUID);
DROP FUNCTION IF EXISTS get_user_chat_reaction(UUID, UUID);

DROP INDEX IF EXISTS idx_chat_likes_chat_id;
DROP INDEX IF EXISTS idx_chat_likes_user_id;
DROP INDEX IF EXISTS idx_chat_likes_reaction_type;
DROP INDEX IF EXISTS idx_chat_likes_created_at;
DROP INDEX IF EXISTS idx_chat_likes_chat_reaction;


DROP TABLE IF EXISTS chats CASCADE;

COMMIT;