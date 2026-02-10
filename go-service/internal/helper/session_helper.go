package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"

	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SessionHelper struct {
	client *redisclient.RedisClient
	logger *zap.Logger
}

func NewSessionHelper(client *redisclient.RedisClient, logger *zap.Logger) *SessionHelper {
	return &SessionHelper{client: client, logger: logger}
}
func (ts *SessionHelper) StoreSessionCSRF(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.SetAny(subCtx, keys, token, 1*time.Minute)
}
func (ts *SessionHelper) GetCSRFBySession(ctx context.Context, token string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.GetString(subCtx, keys)
}
func (ts *SessionHelper) ClearCSRFBySession(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("csrf:%s", token)
	return ts.client.DeleteCache(subCtx, keys)
}
func (ts *SessionHelper) SetSessionToken(ctx context.Context, user *entity.User, expiry time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", user.ID)

	ts.logger.Info("💾 [SET] Storing session in Redis",
		zap.String("key", keys),
		zap.String("user_id", user.ID.String()),
		zap.Duration("expiry", expiry))

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	return ts.client.Setbyte(subCtx, keys, userJSON, expiry)
}
func (ts *SessionHelper) GetSessionTokenMetadata(ctx context.Context, userID string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", userID)

	ts.logger.Info("🔍 [GET] Looking up session in Redis",
		zap.String("key", keys),
		zap.String("session_id", userID))

	data, err := ts.client.GetByte(subCtx, keys)
	if err != nil {
		ts.logger.Error("❌ [GET FAILED] Redis lookup failed",
			zap.String("key", keys),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get session metadata from cache storage: %w", err)
	}

	var user entity.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		ts.logger.Error("❌ [GET FAILED] Unmarshal failed",
			zap.String("key", keys),
			zap.String("raw_data", string(data)),
			zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	ts.logger.Info("✅ [GET SUCCESS] Session found",
		zap.String("key", keys),
		zap.String("user_id", user.ID.String()))

	return &user, nil
}
func (ts *SessionHelper) IsSessionTokenRevoked(ctx context.Context, sessionID string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("revoked_token:%s", sessionID)
	return ts.client.GetBool(subCtx, keys)
}
func (ts *SessionHelper) RevokeSessionToken(ctx context.Context, sessionID string, remainingTTL time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("revoked_token:%s", sessionID)
	return ts.client.SetAny(subCtx, keys, "revoked", remainingTTL)
}
func (ts *SessionHelper) ClearSessionToken(ctx context.Context, sessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	keys := fmt.Sprintf("session:%s", sessionID)
	return ts.client.DeleteCache(subCtx, keys)
}
func (ts *SessionHelper) DeleteAllExceptCurrent(ctx context.Context, currentSessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	currentMeta, err := ts.GetSessionTokenMetadata(subCtx, currentSessionID)
	if err != nil {
		return fmt.Errorf("failed to get current session metadata: %w", err)
	}

	var deletedCount int
	var cursor uint64
	var lastErr error
	pattern := "session:*"

	for {
		keys, nextCursor, err := ts.client.Client().Scan(subCtx, cursor, pattern, 100).Result()
		if err != nil {
			lastErr = fmt.Errorf("scan failed: %w", err)
			break
		}
		for _, key := range keys {
			if key == fmt.Sprintf("session:%s", currentSessionID) {
				continue
			}

			sessionID := strings.TrimPrefix(key, "session:")
			meta, err := ts.GetSessionTokenMetadata(subCtx, sessionID)

			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				ts.logger.Warn("Failed to get session metadata",
					zap.String("key", key),
					zap.Error(err))
				continue
			}

			if meta.ID == currentMeta.ID {
				if err := ts.client.DeleteCache(subCtx, key); err != nil {
					lastErr = fmt.Errorf("failed to delete session %s: %w", key, err)
					continue
				}
				deletedCount++

				revokeKey := fmt.Sprintf("revoked_token:%s", sessionID)
				ts.client.DeleteCache(subCtx, revokeKey)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	ts.logger.Info("Deleted old sessions",
		zap.Int("deleted_count", deletedCount),
		zap.String("current_session", currentSessionID),
		zap.String("user_id", currentMeta.ID.String()))

	return lastErr

}
func (ts *SessionHelper) ClearAllSessionToken(ctx context.Context, sessionID, csrfToken string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	csrfKey := fmt.Sprintf("csrf:%s", csrfToken)
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	revokedKey := fmt.Sprintf("revoked_token:%s", sessionID)

	_, err := ts.client.Client().Pipelined(subCtx, func(pipe redis.Pipeliner) error {
		pipe.Del(subCtx, csrfKey)
		pipe.Del(subCtx, sessionKey)
		pipe.Del(subCtx, revokedKey)
		return nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("redis operation timed out: %w", err)
		}
		return fmt.Errorf("failed to clear session tokens: %w", err)
	}

	return nil
}
func (ts *SessionHelper) UpdateLastAccessed(ctx context.Context, sessionID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	metadata, err := ts.GetSessionTokenMetadata(subCtx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session data: %w", err)
	}

	key := fmt.Sprintf("session:%s", sessionID)
	ttl, err := ts.client.Client().TTL(subCtx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get TTL: %w", err)
	}

	if ttl > 0 {
		return ts.SetSessionToken(ctx, metadata, ttl)
	}

	return nil
}
func (ts *SessionHelper) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	var sessions []*entity.User
	var cursor uint64
	pattern := "session:*"

	for {
		keys, nextCursor, err := ts.client.Client().Scan(subCtx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		for _, key := range keys {
			sessionID := strings.TrimPrefix(key, "session:")
			metadata, err := ts.GetSessionTokenMetadata(subCtx, sessionID)
			if err != nil {
				continue // Skip invalid sessions
			}

			if metadata.ID == userID {
				sessions = append(sessions, metadata)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return sessions, nil
}
func (ts *SessionHelper) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	sessions, err := ts.GetUserSessions(subCtx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := ts.ClearSessionToken(subCtx, session.ID.String()); err != nil {
			ts.logger.Warn("Failed to delete session",
				zap.String("session_id", session.ID.String()),
				zap.Error(err))
		}

		revokeKey := fmt.Sprintf("revoked_token:%s", session.ID.String())
		ts.client.DeleteCache(subCtx, revokeKey)
	}

	ts.logger.Info("Deleted all user sessions",
		zap.String("user_id", userID.String()),
		zap.Int("session_count", len(sessions)))

	return nil
}
