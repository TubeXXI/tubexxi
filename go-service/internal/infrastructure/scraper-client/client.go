package scraper_client

import (
	"context"
	"fmt"
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

func (s *ScraperClient) GetMoviesByGenre(ctx context.Context, slug string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetMoviesByGenre(ctx, &pb.GenreRequest{Slug: slug, Page: page})
}

func (s *ScraperClient) SearchMovies(ctx context.Context, query string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.SearchMovies(ctx, &pb.SearchRequest{Query: query, Page: page})
}

func (s *ScraperClient) GetMoviesByFeature(ctx context.Context, featureType string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetMoviesByFeature(ctx, &pb.FeatureRequest{FeatureType: featureType, Page: page})
}

func (s *ScraperClient) GetMoviesByCountry(ctx context.Context, country string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetMoviesByCountry(ctx, &pb.CountryRequest{CountrySlug: country, Page: page})
}

func (s *ScraperClient) GetMoviesByYear(ctx context.Context, year int32, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetMoviesByYear(ctx, &pb.YearRequest{Year: year, Page: page})
}

func (s *ScraperClient) GetSpecialPage(ctx context.Context, pageName string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSpecialPage(ctx, &pb.SpecialPageRequest{PageName: pageName, Page: page})
}

func (s *ScraperClient) GetMovieDetail(ctx context.Context, slug string) (*pb.MovieDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetMovieDetail(ctx, &pb.MovieDetailRequest{Slug: slug})
}

