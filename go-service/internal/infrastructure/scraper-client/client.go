package scraper_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"tubexxi/video-api/config"
	pb "tubexxi/video-api/proto"
)

type ScraperClient struct {
	client pb.ScraperServiceClient
	conn   *grpc.ClientConn
}

func NewScraperClient(cfg config.ScraperConfig) (*ScraperClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to scraper service: %w", err)
	}

	client := pb.NewScraperServiceClient(conn)
	return &ScraperClient{
		client: client,
		conn:   conn,
	}, nil
}

func (s *ScraperClient) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *ScraperClient) ScrapeHome(ctx context.Context) (*pb.HomeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.client.ScrapeHome(ctx, &pb.Empty{})
}

func (s *ScraperClient) ScrapeList(ctx context.Context, url string) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.client.ScrapeList(ctx, &pb.ListRequest{Url: url})
}
