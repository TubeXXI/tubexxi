package middleware

import (
	"errors"
	"strings"
	"time"

	"tubexxi/video-api/internal/entity"
	firebaseclient "tubexxi/video-api/internal/infrastructure/firebase-client"
	"tubexxi/video-api/internal/infrastructure/metrics"
	"tubexxi/video-api/internal/infrastructure/repository"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/telegram"
	"tubexxi/video-api/pkg/utils"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
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
	notifier     telegram.Notifier
	ctxinject    *ContextMiddleware
	firebase     *firebaseclient.FirebaseClient
	roleRepo     repository.RoleRepository
	userRepo     repository.UserRepository
	logger       *zap.Logger
	checkRevoked bool
}

func NewAuthMiddleware(
	notifier telegram.Notifier,
	ctxinject *ContextMiddleware,
	firebase *firebaseclient.FirebaseClient,
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
	checkRevoked bool,
) *AuthMiddleware {
	return &AuthMiddleware{
		notifier:     notifier,
		ctxinject:    ctxinject,
		firebase:     firebase,
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		logger:       logger,
		checkRevoked: checkRevoked,
	}
}

func (m *AuthMiddleware) FirebaseAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		m.logger.Debug("🔐 [START] Firebase Authentication", zap.String("path", c.Path()))

		tokenStr, err := m.extractTokenFromHeader(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		firebaseToken, err := m.firebase.VerifyIDToken(ctx, tokenStr, m.checkRevoked)
		if err != nil {
			m.logger.Error("❌ [STEP 1 FAILED] Firebase token verification failed",
				zap.Error(err),
				zap.String("path", c.Path()))
			return m.handleTokenError(c, err)
		}

		email, _ := firebaseToken.Claims["email"].(string)
		if email == "" {
			return response.Error(c, fiber.StatusUnauthorized, ErrInvalidToken, nil)
		}

		name, _ := firebaseToken.Claims["name"].(string)
		picture, _ := firebaseToken.Claims["picture"].(string)
		phone, _ := firebaseToken.Claims["phone_number"].(string)
		emailVerified, _ := firebaseToken.Claims["email_verified"].(bool)

		user, err := m.userRepo.FindByEmail(ctx, email)
		if err != nil {
			randomPassword := utils.GenerateRandomPassword()
			hash, hashErr := utils.HashPassword(randomPassword)
			if hashErr != nil {
				m.logger.Error("failed to hash fallback password", zap.Error(hashErr))
				return response.Error(c, fiber.StatusInternalServerError, "internal_error", nil)
			}

			nm := name
			if nm == "" {
				nm = "User"
			}

			newUser := &entity.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: hash,
				FullName:     nm,
				Phone:        entity.NewNullString(phone),
				AvatarURL:    entity.NewNullString(picture),
				IsActive:     true,
				IsVerified:   emailVerified,
				CreatedAt:    time.Now(),
			}

			role, err := m.roleRepo.FindByLevel(ctx, entity.RoleLevelUser)
			if err != nil {
				m.logger.Error("failed to find role: %w", zap.Error(err))
				return response.Error(c, fiber.StatusInternalServerError, "internal_error", nil)
			}
			newUser.RoleID = role.ID

			if createErr := m.userRepo.CreateWithRecovery(ctx, newUser); createErr != nil {
				m.logger.Error("failed to create user from firebase token", zap.Error(createErr), zap.String("email", email))
				return response.Error(c, fiber.StatusUnauthorized, ErrInvalidToken, nil)
			}
			user = newUser
		}

		if emailVerified && !user.IsVerified {
			_ = m.userRepo.SetEmailVerified(ctx, user.ID, true)
			user.IsVerified = true
		}

		changed := false
		if name != "" && (strings.TrimSpace(user.FullName) == "" || user.FullName == "User") {
			user.FullName = name
			changed = true
		}
		if phone != "" && !user.Phone.Valid {
			user.Phone = entity.NewNullString(phone)
			changed = true
		}
		if picture != "" && !user.AvatarURL.Valid {
			user.AvatarURL = entity.NewNullString(picture)
			changed = true
		}
		if changed {
			if updated, upErr := m.userRepo.UpdateWithRecovery(ctx, user); upErr == nil && updated != nil {
				user = updated
			}
		}

		var roleLevel int
		var roleName string
		if user.Role != nil && user.Role.ID != uuid.Nil {
			roleLevel = int(user.Role.Level)
			roleName = user.Role.Name
		} else {
			roleLevel, roleName = extractRoleFromClaims(firebaseToken.Claims)
		}

		c.Locals("user_id", user.ID.String())
		c.Locals("email", user.Email)
		c.Locals("role_level", roleLevel)
		c.Locals("role_name", roleName)
		c.Locals("firebase_uid", firebaseToken.UID)
		c.Locals("firebase_claims", firebaseToken.Claims)

		m.logger.Debug("✅ [FINAL] Firebase authentication successful",
			zap.String("user_id", user.ID.String()),
			zap.Int("role_level", roleLevel),
			zap.String("role_name", roleName),
			zap.String("firebase_uid", firebaseToken.UID),
			zap.String("path", c.Path()))

		return c.Next()
	}
}

func extractRoleFromClaims(claims map[string]interface{}) (int, string) {
	if claims == nil {
		return entity.RoleLevelUser, entity.RoleUser
	}
	if v, ok := claims["role_level"]; ok {
		switch t := v.(type) {
		case float64:
			lv := int(t)
			return lv, roleNameFromLevel(lv)
		case int:
			return t, roleNameFromLevel(t)
		}
	}
	if v, ok := claims["admin"].(bool); ok && v {
		return entity.RoleLevelAdmin, entity.RoleAdmin
	}
	if v, ok := claims["role"].(string); ok {
		s := strings.ToLower(strings.TrimSpace(v))
		switch s {
		case entity.RoleSuperAdmin:
			return entity.RoleLevelSuperAdmin, entity.RoleSuperAdmin
		case entity.RoleAdmin:
			return entity.RoleLevelAdmin, entity.RoleAdmin
		default:
			return entity.RoleLevelUser, entity.RoleUser
		}
	}
	return entity.RoleLevelUser, entity.RoleUser
}

func roleNameFromLevel(level int) string {
	switch level {
	case entity.RoleLevelSuperAdmin:
		return entity.RoleSuperAdmin
	case entity.RoleLevelAdmin:
		return entity.RoleAdmin
	default:
		return entity.RoleUser
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
func (m *AuthMiddleware) handleTokenError(c *fiber.Ctx, err error) error {
	errorMsg := ErrInvalidToken

	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("firebase_validation_failed").Inc()

	return response.Error(c, fiber.StatusUnauthorized, errorMsg, nil)
}
func (am *AuthMiddleware) LogUnauthorized(c *fiber.Ctx, subject string, requestID string) {
	metrics.GetAppMetrics().JWTErrorTotal.WithLabelValues("firebase_invalid_token").Inc()
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
