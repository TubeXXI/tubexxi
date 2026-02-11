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

// Series methods

func (s *ScraperClient) GetSeriesHome(ctx context.Context) (*pb.HomeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesHome(ctx, &pb.Empty{})
}

func (s *ScraperClient) GetSeriesByGenre(ctx context.Context, slug string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesByGenre(ctx, &pb.GenreRequest{Slug: slug, Page: page})
}

func (s *ScraperClient) SearchSeries(ctx context.Context, query string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.SearchSeries(ctx, &pb.SearchRequest{Query: query, Page: page})
}

func (s *ScraperClient) GetSeriesByFeature(ctx context.Context, featureType string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesByFeature(ctx, &pb.FeatureRequest{FeatureType: featureType, Page: page})
}

func (s *ScraperClient) GetSeriesByCountry(ctx context.Context, country string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesByCountry(ctx, &pb.CountryRequest{CountrySlug: country, Page: page})
}

func (s *ScraperClient) GetSeriesByYear(ctx context.Context, year int32, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesByYear(ctx, &pb.YearRequest{Year: year, Page: page})
}

func (s *ScraperClient) GetSeriesSpecialPage(ctx context.Context, pageName string, page int32) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesSpecialPage(ctx, &pb.SpecialPageRequest{PageName: pageName, Page: page})
}

func (s *ScraperClient) GetSeriesDetail(ctx context.Context, slug string) (*pb.SeriesDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesDetail(ctx, &pb.MovieDetailRequest{Slug: slug})
}
func (s *ScraperClient) GetSeriesEpisode(ctx context.Context, url string) (*pb.SeriesEpisodeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetSeriesEpisode(ctx, &pb.SeriesEpisodeRequest{Url: url})
}

// Anime methods

func (s *ScraperClient) GetAnimeLatest(ctx context.Context, page int32) (*pb.AnimeListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetAnimeLatest(ctx, &pb.AnimePageRequest{Page: page})
}

func (s *ScraperClient) SearchAnime(ctx context.Context, query string, page int32) (*pb.AnimeListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.SearchAnime(ctx, &pb.AnimeSearchRequest{Query: query, Page: page})
}

func (s *ScraperClient) GetAnimeOngoing(ctx context.Context, page int32) (*pb.AnimeListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetAnimeOngoing(ctx, &pb.AnimePageRequest{Page: page})
}

func (s *ScraperClient) GetAnimeGenres(ctx context.Context) (*pb.AnimeGenresResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetAnimeGenres(ctx, &pb.Empty{})
}

func (s *ScraperClient) GetAnimeDetail(ctx context.Context, url string) (*pb.AnimeDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetAnimeDetail(ctx, &pb.AnimeUrlRequest{Url: url})
}

func (s *ScraperClient) GetAnimeEpisode(ctx context.Context, url string) (*pb.AnimeEpisodeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.GetAnimeEpisode(ctx, &pb.AnimeUrlRequest{Url: url})
}
