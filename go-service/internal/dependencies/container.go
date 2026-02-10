package dependencies

import (
	"log"

	"tubexxi/video-api/config"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	"tubexxi/video-api/internal/infrastructure/metrics"
	"tubexxi/video-api/internal/infrastructure/redis-client"
	scraper_client "tubexxi/video-api/internal/infrastructure/scraper-client"
)

type Container struct {
	Config        *config.Config
	Redis         *redis_client.RedisClient
	ContextPool   *contextpool.ContextPool
	Metrics       metrics.MetricService
	ScraperClient *scraper_client.ScraperClient
}

func NewContainer() *Container {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient, err := redis_client.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize Context Pool
	contextPool := contextpool.NewContextPool()

	// Initialize Metrics
	metricService := metrics.NewMetricFactory().Create("app")
	
	// Initialize Scraper Client
	scraperClient, err := scraper_client.NewScraperClient(cfg.Scraper)
	if err != nil {
		log.Printf("Warning: Failed to initialize Scraper Client: %v", err)
	}

	return &Container{
		Config:        cfg,
		Redis:         redisClient,
		ContextPool:   contextPool,
		Metrics:       metricService,
		ScraperClient: scraperClient,
	}
}
