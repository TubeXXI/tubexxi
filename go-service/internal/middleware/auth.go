package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tubexxi/video-api/internal/entity"
	helpers "tubexxi/video-api/internal/helper"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	"tubexxi/video-api/internal/infrastructure/metrics"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/telegram"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrMissingToken   = "Unauthorized - Missing token"
	ErrMissMatchToken = "Unauthorized - Token mismatch"
	ErrInvalidToken   = "Unauthorized - Invalid token format"
	ErrExpiredToken   = "Unauthorized - Invalid or expired token"
	ErrInvalidAuth    = "Unauthorized - Invalid authorization header"
)

type AuthMiddleware struct {
	notifier      telegram.Notifier
	ctxinject     *ContextMiddleware
	redisClient   *redisclient.RedisClient
	sessionHelper *helpers.SessionHelper
	logger        *zap.Logger
	jwtSecret     string
}

func NewAuthMiddleware(notifier telegram.Notifier, ctxinject *ContextMiddleware, redisClient *redisclient.RedisClient, sessionHelper *helpers.SessionHelper, logger *zap.Logger, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		notifier:      notifier,
		ctxinject:     ctxinject,
		redisClient:   redisClient,
		sessionHelper: sessionHelper,
		logger:        logger,
		jwtSecret:     jwtSecret,
	}
}

func (m *AuthMiddleware) JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		m.logger.Debug("🔐 [START] JWT Authentication", zap.String("path", c.Path()))

		tokenStr, err := m.extractTokenFromHeader(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		m.logger.Debug("🔐 [STEP 1] Extracted token from header",
			zap.String("token_prefix", tokenStr[:20]+"..."),
			zap.String("path", c.Path()))

		jwtClaims, sessionID, err := m.validateJWTAndExtractSession(tokenStr)
		if err != nil {
			m.logger.Error("❌ [STEP 1 FAILED] JWT validation failed",
				zap.Error(err),
				zap.String("token_prefix", tokenStr[:20]+"..."))
			return m.handleTokenError(c, err)
		}
		m.logger.Debug("✅ [STEP 1 SUCCESS] JWT validated",
			zap.String("session_id", sessionID),
			zap.String("user_id", jwtClaims["sub"].(string)))

		userData, err := m.validateTokenSession(ctx, sessionID)
		if err != nil {
			m.logger.Error("❌ [STEP 2 FAILED] Redis session validation failed",
				zap.Error(err),
				zap.String("session_id", sessionID),
				zap.String("redis_key", "session:"+sessionID))
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		m.logger.Debug("✅ [STEP 2 SUCCESS] Redis session found",
			zap.String("user_id", userData.ID.String()),
			zap.String("email", userData.Email))

		if err := m.crossValidateClaims(jwtClaims, userData); err != nil {
			m.logger.Warn("❌ [STEP 3 FAILED] Token claim mismatch",
				zap.String("user_id", userData.ID.String()),
				zap.Error(err))
			return response.Error(c, fiber.StatusUnauthorized, ErrInvalidToken, nil)
		}

		m.logger.Debug("✅ [STEP 3 SUCCESS] Claims cross-validated")

		m.setContextLocals(c, userData)

		m.logger.Debug("✅ [FINAL] JWT authentication successful",
			zap.String("user_id", userData.ID.String()),
			zap.String("session_id", sessionID),
			zap.String("path", c.Path()))

		return c.Next()
	}
}
func (m *AuthMiddleware) extractTokenFromHeader(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(ErrMissingToken)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New(ErrInvalidAuth)
	}

	if parts[1] == "" {
		return "", errors.New(ErrMissingToken)
	}

	return parts[1], nil
}
func (m *AuthMiddleware) validateJWTAndExtractSession(tokenStr string) (jwt.MapClaims, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, "", err
	}

	if !token.Valid {
		return nil, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", errors.New("invalid token claims")
	}

	if exp, okExp := claims["exp"].(float64); okExp {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, "", jwt.ErrTokenExpired
		}
	} else {
		return nil, "", errors.New("missing expiration claim")
	}

	if _, okSub := claims["sub"].(string); !okSub {
		return nil, "", errors.New("missing subject claim")
	}

	sessionID, okSid := claims["sid"].(string)
	if !okSid || sessionID == "" {
		return nil, "", errors.New("missing session ID in token")
	}

	return claims, sessionID, nil
}
func (m *AuthMiddleware) validateTokenSession(ctx context.Context, userID string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	m.logger.Debug("🔍 [validateTokenSession] Looking up session",
		zap.String("user_id", userID))

	userSession, err := m.sessionHelper.GetSessionTokenMetadata(subCtx, userID)
	if err != nil {
		m.logger.Error("❌ [validateTokenSession] GetSessionTokenMetadata failed",
			zap.Error(err),
			zap.String("user_id", userID))
		return nil, errors.New(ErrInvalidToken)
	}

	m.logger.Debug("🔍 [validateTokenSession] Redis data retrieved",
		zap.String("user_id", userSession.ID.String()),
		zap.String("email", userSession.Email))

	if err := m.validateTokenMetadata(userSession); err != nil {
		m.logger.Error("❌ [validateTokenSession] validateTokenMetadata failed",
			zap.Error(err),
			zap.String("user_id", userSession.ID.String()))
		return nil, err
	}

	return userSession, nil
}

func (m *AuthMiddleware) handleTokenError(c *fiber.Ctx, err error) error {
	errorMsg := ErrInvalidToken
	if errors.Is(err, jwt.ErrTokenExpired) {
		errorMsg = ErrExpiredToken
	}

	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("validation_failed").Inc()

	return response.Error(c, fiber.StatusUnauthorized, errorMsg, nil)
}
func (m *AuthMiddleware) crossValidateClaims(jwtClaims jwt.MapClaims, userSession *entity.User) error {
	jwtUserIDStr, ok := jwtClaims["sub"].(string)
	if !ok {
		return errors.New("missing user ID in JWT claims")
	}

	jwtUserID, err := uuid.Parse(jwtUserIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID format in JWT: %w", err)
	}

	if jwtUserID != userSession.ID {
		return fmt.Errorf("user ID mismatch: JWT=%s, Redis=%s", jwtUserID, userSession.ID)
	}

	if jwtSessionID, ok := jwtClaims["sid"].(string); ok && jwtSessionID != "" {
		if jwtSessionID != userSession.ID.String() {
			return fmt.Errorf("session ID mismatch: JWT=%s, Redis=%s", jwtSessionID, userSession.ID.String())
		}
	} else {
		return errors.New("missing session ID in JWT claims")
	}

	if jwtRoleIDStr, ok := jwtClaims["rid"].(string); ok && jwtRoleIDStr != "" {
		jwtRoleID, err := uuid.Parse(jwtRoleIDStr)
		if err != nil {
			return fmt.Errorf("invalid role ID in JWT: %w", err)
		}
		if jwtRoleID != userSession.Role.ID {
			return fmt.Errorf("role ID mismatch: JWT=%s, Redis=%s", jwtRoleID, userSession.Role.ID)
		}
	} else {
		return errors.New("missing role ID in JWT claims")
	}

	if jwtEmail, ok := jwtClaims["em"].(string); ok && jwtEmail != "" {
		if jwtEmail != userSession.Email {
			return fmt.Errorf("email mismatch: JWT=%s, Redis=%s", jwtEmail, userSession.Email)
		}
	}

	return nil
}
func (m *AuthMiddleware) validateTokenMetadata(metadata *entity.User) error {
	if metadata.ID == uuid.Nil {
		return fmt.Errorf("invalid user ID in session")
	}
	if metadata.Role.ID == uuid.Nil {
		return fmt.Errorf("invalid role ID in session")
	}

	return nil
}
func (m *AuthMiddleware) setContextLocals(c *fiber.Ctx, metadata *entity.User) {
	c.Locals("user_id", metadata.ID.String())
	c.Locals("role_level", metadata.Role.Level)
	c.Locals("email", metadata.Email)

	if len(metadata.Role.Name) > 0 {
		c.Locals("role_name", metadata.Role.Name)
	}
}
func (am *AuthMiddleware) LogUnauthorized(c *fiber.Ctx, subject string, requestID string) {
	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("invalid_signature").Inc()
	am.notifier.SendAlert(telegram.AlertRequest{
		Subject: subject,
		Message: subject,
		Metadata: map[string]interface{}{
			"request_id": requestID,
			"timestamp":  time.Now(),
			"user_agent": c.Get("User-Agent"),
			"ip":         c.Locals("real_ip").(string),
		},
	})
}
