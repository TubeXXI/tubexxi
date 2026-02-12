package middleware

import (
	"fmt"
	"time"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

type RateLimiterMiddleware struct {
	ctxinjext   *ContextMiddleware
	redisClient *redisclient.RedisClient
	logger      *zap.Logger
}

func NewRateLimiterMiddleware(
	ctxinject *ContextMiddleware,
	redisClient *redisclient.RedisClient,
	logger *zap.Logger,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		ctxinjext:   ctxinject,
		redisClient: redisClient,
		logger:      logger,
	}
}
func (rm *RateLimiterMiddleware) ProgressDelay(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := rm.ctxinjext.From(c)
		rm.logger.Debug("Progressing delay", zap.String("key", key), zap.String("path", c.Path()))

		attemptsKey := fmt.Sprintf("delay:%s:%s", key, c.Locals("real_ip").(string))

		attempts, err := rm.redisClient.Increment(ctx, attemptsKey)
		if err != nil {
			rm.logger.Error("Redis error", zap.Error(err))
			return c.Next()
		}
		_ = rm.redisClient.Expire(ctx, attemptsKey, 30*time.Minute)

		if attempts >= 4 {
			delay := time.Duration(attempts-3) * time.Second
			time.Sleep(delay)
		}
		return c.Next()
	}
}
func (rm *RateLimiterMiddleware) ResetLimitCounters(c *fiber.Ctx) {
	ctx := rm.ctxinjext.From(c)
	rm.logger.Debug("Resetting limit counters", zap.String("path", c.Path()))

	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	ip := c.Locals("real_ip").(string)
	patterns := []string{
		fmt.Sprintf("rate:%s:%s", "login", ip),
		fmt.Sprintf("delay:%s:%s", "login", ip),
		fmt.Sprintf("rate:%s:%s", "register", ip),
		fmt.Sprintf("delay:%s:%s", "register", ip),
		fmt.Sprintf("rate:%s:%s", "reset_password", ip),
		fmt.Sprintf("delay:%s:%s", "reset_password", ip),
		fmt.Sprintf("rate:%s:%s", "verify_email", ip),
		fmt.Sprintf("block:%s:%s", "change_password", ip),
	}

	for _, pattern := range patterns {
		rm.redisClient.DeleteCache(ctx, pattern)
	}
}
func (rm *RateLimiterMiddleware) GlobalRequestLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("global:%s:%s", ip, c.Method())
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		Storage:    rm.redisClient,
		LimitReached: func(c *fiber.Ctx) error {
			retryAfter := c.GetRespHeader("Retry-After")
			return response.Error(c, fiber.StatusTooManyRequests, "global_rate_limit_exceeded", fiber.Map{
				"retry_after": retryAfter,
			})
		},
	})
}
func (rm *RateLimiterMiddleware) BaseLimiter(key string, max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("rate:%s:%s", key, ip)
		},
		Storage:      rm.redisClient,
		Max:          max,
		Expiration:   expiration,
		LimitReached: defaultLimitReachedHandler,
	})
}
func (rm *RateLimiterMiddleware) BlockLimiter(key string, maxAttempts int, blockDuration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Locals("real_ip").(string)
			return fmt.Sprintf("block:%s:%s", key, ip)
		},
		Storage:      rm.redisClient,
		Max:          maxAttempts,
		Expiration:   blockDuration,
		LimitReached: blockLimitReachedHandler(blockDuration),
	})
}
func defaultLimitReachedHandler(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusTooManyRequests, "You have reached the request limit. Please try again later.", nil)
}
func blockLimitReachedHandler(blockDuration time.Duration) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Too many attempts. Please try again after %v.", blockDuration)
		return response.Error(c, fiber.StatusTooManyRequests, msg, fiber.Map{
			"block_duration": blockDuration,
		})
	}
}
