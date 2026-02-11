package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"tubexxi/video-api/internal/entity"
	scraper_client "tubexxi/video-api/internal/infrastructure/scraper-client"
	pb "tubexxi/video-api/proto"
)

type MovieService struct {
	logger        *zap.Logger
	scraperClient *scraper_client.ScraperClient
}

func NewMovieService(
	logger *zap.Logger,
	scraperClient *scraper_client.ScraperClient,
) *MovieService {
	return &MovieService{
		logger:        logger,
		scraperClient: scraperClient,
	}
}

func (s *MovieService) GetHome(ctx context.Context) ([]entity.HomeScrapperResponse, error) {
	// Call Python Scraper Service
	resp, err := s.scraperClient.ScrapeHome(ctx)
	if err != nil {
		s.logger.Error("failed to scrape home", zap.Error(err))
		return nil, err
	}

	// Map Proto response to Entity
	var homeResponse []entity.HomeScrapperResponse
	for _, section := range resp.Sections {
		var movies []entity.Movie
		for _, m := range section.Value {
			movie := s.mapProtoToMovie(m)
			movies = append(movies, movie)
		}

		// Handle optional ViewAllUrl
		var viewAllUrl *string
		if section.ViewAllUrl != nil {
			viewAllUrl = section.ViewAllUrl
		}

		homeResponse = append(homeResponse, entity.HomeScrapperResponse{
			Key:        section.Key,
			Value:      movies,
			ViewAllUrl: viewAllUrl,
		})
	}

	return homeResponse, nil
}

func (s *MovieService) GetList(ctx context.Context, url string) (*entity.HomeScrapperListResponse, error) {
	// Placeholder implementation for List
	// We need to define a proper List entity response that includes pagination
	// For now, let's just return nil to satisfy interface if we had one, or implement later
	return nil, fmt.Errorf("not implemented")
}

func (s *MovieService) mapProtoToMovie(m *pb.Movie) entity.Movie {
	// Helper to map pointer fields

	// UUID Parsing
	id, err := uuid.Parse(m.Id)
	if err != nil {
		id = uuid.New() // Fallback if invalid
	}

	// Rating (float32 -> float64 pointer)
	var rating *float64
	if m.Rating != nil {
		r := float64(*m.Rating)
		rating = &r
	}

	// Duration (int64 -> int64 pointer)
	var duration *int64
	if m.Duration != nil {
		d := *m.Duration
		duration = &d
	}

	// Year (int32 -> int32 pointer)
	var year *int32
	if m.Year != nil {
		y := *m.Year
		year = &y
	}

	return entity.Movie{
		ID:              id,
		Title:           m.Title,
		OriginalTitle:   m.OriginalTitle,
		Thumbnail:       m.Thumbnail,
		Synopsis:        m.Synopsis,
		Rating:          rating,
		Duration:        duration,
		Year:            year,
		DatePublished:   nil, // Skip for now as parsing is complex
		LabelQuality:    m.LabelQuality,
		Genre:           m.Genre,
		OriginalPageUrl: m.OriginalPageUrl,
	}
}
