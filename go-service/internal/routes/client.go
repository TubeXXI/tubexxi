package routes

import (
	"time"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type ClientRoutes struct {
	path    string
	handler *handler.ClientHandler
	limiter *middleware.RateLimiterMiddleware
	csrf    *middleware.CSRFMiddleware
}

func NewClientRoutes(
	handler *handler.ClientHandler,
	limiter *middleware.RateLimiterMiddleware,
	csrf *middleware.CSRFMiddleware,
) *ClientRoutes {
	return &ClientRoutes{
		path:    "/client",
		handler: handler,
		limiter: limiter,
		csrf:    csrf,
	}
}
func (r *ClientRoutes) RegisterRoutes(parent fiber.Router) {
	client := parent.Group(r.path)

	webClient := client.Group("/web")

	webClient.Post("/contact",
		r.limiter.BlockLimiter("send_contact",
			3,
			5*time.Minute,
		),
		r.csrf.CSRFProtect(),
		r.handler.SendContactEmail,
	)
}
