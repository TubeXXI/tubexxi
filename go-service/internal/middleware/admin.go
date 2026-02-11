package middleware

import (
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AdminMiddleware struct {
	ctxinject *ContextMiddleware
	logger    *zap.Logger
}

func NewAdminMiddleware(ctxinject *ContextMiddleware, logger *zap.Logger) *AdminMiddleware {
	return &AdminMiddleware{
		logger:    logger,
		ctxinject: ctxinject,
	}
}

func (m *AdminMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
		}

		firebaseUID, _ := c.Locals("firebase_uid").(string)
		if firebaseUID == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Firebase authentication required", nil)
		}

		role, ok := c.Locals("role_level").(int)
		if !ok || (role != entity.RoleLevelAdmin && role != entity.RoleLevelSuperAdmin) {
			return response.Error(c, fiber.StatusForbidden, "Admin access required", nil)
		}

		return c.Next()

	}
}
