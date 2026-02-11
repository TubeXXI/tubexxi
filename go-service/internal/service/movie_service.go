package service

import (
	"context"
	"time"

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
	resp, err := s.scraperClient.ScrapeHome(ctx)
	if err != nil {
		s.logger.Error("failed to scrape home", zap.Error(err))
		return nil, err
	}

	var homeResponse []entity.HomeScrapperResponse
	for _, section := range resp.Sections {
		var movies []entity.Movie
		for _, m := range section.Value {
			movie := s.mapProtoToMovie(m)
			movies = append(movies, movie)
		}

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

func (s *MovieService) GetMoviesByGenre(ctx context.Context, slug string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetMoviesByGenre(ctx, slug, page)
	if err != nil {
		s.logger.Error("failed to get movies by genre", zap.String("slug", slug), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) SearchMovies(ctx context.Context, query string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.SearchMovies(ctx, query, page)
	if err != nil {
		s.logger.Error("failed to search movies", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetMoviesByFeature(ctx context.Context, featureType string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetMoviesByFeature(ctx, featureType, page)
	if err != nil {
		s.logger.Error("failed to get movies by feature", zap.String("type", featureType), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetMoviesByCountry(ctx context.Context, country string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetMoviesByCountry(ctx, country, page)
	if err != nil {
		s.logger.Error("failed to get movies by country", zap.String("country", country), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetMoviesByYear(ctx context.Context, year int32, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetMoviesByYear(ctx, year, page)
	if err != nil {
		s.logger.Error("failed to get movies by year", zap.Int32("year", year), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetSpecialPage(ctx context.Context, pageName string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSpecialPage(ctx, pageName, page)
	if err != nil {
		s.logger.Error("failed to get special page", zap.String("page", pageName), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetMovieDetail(ctx context.Context, slug string) (*entity.MovieDetail, error) {
	resp, err := s.scraperClient.GetMovieDetail(ctx, slug)
	if err != nil {
		s.logger.Error("failed to get movie detail", zap.String("slug", slug), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToMovieDetail(resp), nil
}

// Series methods

func (s *MovieService) GetSeriesHome(ctx context.Context) ([]entity.HomeScrapperResponse, error) {
	resp, err := s.scraperClient.GetSeriesHome(ctx)
	if err != nil {
		s.logger.Error("failed to scrape series home", zap.Error(err))
		return nil, err
	}

	var homeResponse []entity.HomeScrapperResponse
	for _, section := range resp.Sections {
		var movies []entity.Movie
		for _, m := range section.Value {
			movie := s.mapProtoToMovie(m)
			movies = append(movies, movie)
		}

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

func (s *MovieService) GetSeriesByGenre(ctx context.Context, slug string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSeriesByGenre(ctx, slug, page)
	if err != nil {
		s.logger.Error("failed to get series by genre", zap.String("slug", slug), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) SearchSeries(ctx context.Context, query string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.SearchSeries(ctx, query, page)
	if err != nil {
		s.logger.Error("failed to search series", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetSeriesByFeature(ctx context.Context, featureType string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSeriesByFeature(ctx, featureType, page)
	if err != nil {
		s.logger.Error("failed to get series by feature", zap.String("type", featureType), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetSeriesByCountry(ctx context.Context, country string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSeriesByCountry(ctx, country, page)
	if err != nil {
		s.logger.Error("failed to get series by country", zap.String("country", country), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetSeriesByYear(ctx context.Context, year int32, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSeriesByYear(ctx, year, page)
	if err != nil {
		s.logger.Error("failed to get series by year", zap.Int32("year", year), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}

func (s *MovieService) GetSeriesSpecialPage(ctx context.Context, pageName string, page int32) (*entity.MovieListResponse, error) {
	resp, err := s.scraperClient.GetSeriesSpecialPage(ctx, pageName, page)
	if err != nil {
		s.logger.Error("failed to get series special page", zap.String("page", pageName), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToListResponse(resp), nil
}


func (s *MovieService) mapProtoToMovie(m *pb.Movie) entity.Movie {
	id, err := uuid.Parse(m.Id)
	if err != nil {
		id = uuid.New()
	}

	var rating *float64
	if m.Rating != nil {
		r := float64(*m.Rating)
		rating = &r
	}

	var duration *int64
	if m.Duration != nil {
		d := *m.Duration
		duration = &d
	}

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
		DatePublished:   nil,
		LabelQuality:    m.LabelQuality,
		Genre:           m.Genre,
		OriginalPageUrl: m.OriginalPageUrl,
	}
}

func (s *MovieService) mapProtoToListResponse(resp *pb.ListResponse) *entity.MovieListResponse {
	var movies []entity.Movie
	for _, m := range resp.Movies {
		movies = append(movies, s.mapProtoToMovie(m))
	}

	var pagination entity.Pagination
	if resp.Pagination != nil {
		pagination = entity.Pagination{
			CurrentPage: resp.Pagination.CurrentPage,
			TotalPage:   resp.Pagination.TotalPage,
			HasNext:     resp.Pagination.HasNext,
			HasPrev:     resp.Pagination.HasPrev,
			NextPageUrl: resp.Pagination.NextPageUrl,
			PrevPageUrl: resp.Pagination.PrevPageUrl,
		}
	}

	return &entity.MovieListResponse{
		Movies:     movies,
		Pagination: pagination,
	}
}

func (s *MovieService) mapProtoToMovieDetail(resp *pb.MovieDetailResponse) *entity.MovieDetail {
	if resp.Detail == nil {
		return nil
	}

	detail := resp.Detail
	movie := s.mapProtoToMovie(detail.Movie)

	var votes *int64
	if detail.Votes != 0 {
		v := detail.Votes
		votes = &v
	}

	var releaseDate *time.Time
	if detail.ReleaseDate != "" {
		if t, err := time.Parse(time.RFC3339, detail.ReleaseDate); err == nil {
			releaseDate = &t
		}
	}

	var updatedAt *time.Time
	if detail.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, detail.UpdatedAt); err == nil {
			updatedAt = &t
		}
	}
	
	// Map trailer url
	var trailerUrl *string
	if detail.TrailerUrl != "" {
		u := detail.TrailerUrl
		trailerUrl = &u
	}

	var playerUrls []entity.PlayerUrl

	for _, p := range detail.PlayerUrls {
		u := p.Url
		t := p.Type
		playerUrls = append(playerUrls, entity.PlayerUrl{
			URL:  &u,
			Type: &t,
		})
	}

	var directors []entity.MoviePerson
	for _, d := range detail.Directors {
		n := d.Name
		pu := d.PageUrl
		directors = append(directors, entity.MoviePerson{
			Name:    &n,
			PageUrl: &pu,
		})
	}

	var stars []entity.MoviePerson
	for _, st := range detail.MovieStars {
		n := st.Name
		pu := st.PageUrl
		stars = append(stars, entity.MoviePerson{
			Name:    &n,
			PageUrl: &pu,
		})
	}

	var countries []entity.CountryMovie
	for _, c := range detail.Countries {
		n := c.Name
		pu := c.PageUrl
		countries = append(countries, entity.CountryMovie{
			Name:    &n,
			PageUrl: &pu,
		})
	}

	var genres []entity.Genre
	for _, g := range detail.Genres {
		n := g.Name
		pu := g.PageUrl
		genres = append(genres, entity.Genre{
			Name:    &n,
			PageUrl: &pu,
		})
	}

	var similarMovies []entity.Movie
	for _, sm := range detail.SimilarMovies {
		similarMovies = append(similarMovies, s.mapProtoToMovie(sm))
	}

	return &entity.MovieDetail{
		Movie:         movie,
		Votes:         votes,
		ReleaseDate:   releaseDate,
		UpdatedAt:     updatedAt,
		TrailerUrl:    trailerUrl,
		PlayerUrl:     &playerUrls,
		Director:      &directors,
		MovieStar:     &stars,
		Countries:     &countries,
		Genres:        &genres,
		SimilarMovies: &similarMovies,
	}
}



