package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"tubexxi/video-api/internal/entity"
	scraper_client "tubexxi/video-api/internal/infrastructure/scraper-client"
	pb "tubexxi/video-api/proto"
)

type AnimeService struct {
	logger       *zap.Logger
	scraperClient *scraper_client.ScraperClient
}

func NewAnimeService(logger *zap.Logger, scraperClient *scraper_client.ScraperClient) *AnimeService {
	return &AnimeService{logger: logger, scraperClient: scraperClient}
}

func (s *AnimeService) GetLatest(ctx context.Context, page int32) (*entity.ScrapeResult, error) {
	resp, err := s.scraperClient.GetAnimeLatest(ctx, page)
	if err != nil {
		s.logger.Error("failed to get anime latest", zap.Int32("page", page), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToScrapeResult(resp), nil
}

func (s *AnimeService) Search(ctx context.Context, query string, page int32) (*entity.ScrapeResult, error) {
	resp, err := s.scraperClient.SearchAnime(ctx, query, page)
	if err != nil {
		s.logger.Error("failed to search anime", zap.String("query", query), zap.Int32("page", page), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToScrapeResult(resp), nil
}

func (s *AnimeService) GetOngoing(ctx context.Context, page int32) (*entity.ScrapeResult, error) {
	resp, err := s.scraperClient.GetAnimeOngoing(ctx, page)
	if err != nil {
		s.logger.Error("failed to get anime ongoing", zap.Int32("page", page), zap.Error(err))
		return nil, err
	}
	return s.mapProtoToScrapeResult(resp), nil
}

func (s *AnimeService) GetGenres(ctx context.Context) (*[]entity.AnimeGenre, error) {
	resp, err := s.scraperClient.GetAnimeGenres(ctx)
	if err != nil {
		s.logger.Error("failed to get anime genres", zap.Error(err))
		return nil, err
	}

	var genres []entity.AnimeGenre
	for _, g := range resp.Genres {
		if g == nil {
			continue
		}
		name := g.Name
		url := g.Url
		var namePtr *string
		if name != "" {
			namePtr = &name
		}
		var urlPtr *string
		if url != "" {
			urlPtr = &url
		}
		genres = append(genres, entity.AnimeGenre{Name: namePtr, URL: urlPtr})
	}

	return &genres, nil
}

func (s *AnimeService) GetDetail(ctx context.Context, url string) (*entity.Anime, error) {
	resp, err := s.scraperClient.GetAnimeDetail(ctx, url)
	if err != nil {
		s.logger.Error("failed to get anime detail", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	if resp.Anime == nil {
		return nil, nil
	}

	anime := s.mapProtoToAnime(resp.Anime)
	return &anime, nil
}

func (s *AnimeService) GetEpisode(ctx context.Context, url string) (*entity.Episode, error) {
	resp, err := s.scraperClient.GetAnimeEpisode(ctx, url)
	if err != nil {
		s.logger.Error("failed to get anime episode", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	if resp.Episode == nil {
		return nil, nil
	}

	ep := s.mapProtoToEpisode(resp.Episode)
	return &ep, nil
}

func (s *AnimeService) mapProtoToScrapeResult(resp *pb.AnimeListResponse) *entity.ScrapeResult {
	if resp == nil {
		return nil
	}

	var animes []entity.Anime
	for _, a := range resp.Animes {
		if a == nil {
			continue
		}
		animes = append(animes, s.mapProtoToAnime(a))
	}

	pagination := entity.PaginationAnime{}
	if resp.Pagination != nil {
		pagination.TotalPages = int(resp.Pagination.TotalPages)
		pagination.CurrentPage = int(resp.Pagination.CurrentPage)
		pagination.HasNext = resp.Pagination.HasNext
		pagination.HasPrevious = resp.Pagination.HasPrevious
		if resp.Pagination.NextPageUrl != "" {
			u := resp.Pagination.NextPageUrl
			pagination.NextPageURL = &u
		}
		if resp.Pagination.PreviousPageUrl != "" {
			u := resp.Pagination.PreviousPageUrl
			pagination.PreviousPageURL = &u
		}
		for _, n := range resp.Pagination.PageNumbers {
			pagination.PageNumbers = append(pagination.PageNumbers, int(n))
		}
		pagination.PerPage = int(resp.Pagination.PerPage)
	}

	var query *string
	if resp.Query != "" {
		q := resp.Query
		query = &q
	}

	return &entity.ScrapeResult{Animes: animes, Pagination: pagination, Query: query}
}

func (s *AnimeService) mapProtoToAnime(a *pb.Anime) entity.Anime {
	setStr := func(v string) *string {
		if v == "" {
			return nil
		}
		vv := v
		return &vv
	}

	var genres []entity.AnimeGenre
	for _, g := range a.Genre {
		if g == nil {
			continue
		}
		name := setStr(g.Name)
		url := setStr(g.Url)
		genres = append(genres, entity.AnimeGenre{Name: name, URL: url})
	}
	var genrePtr *[]entity.AnimeGenre
	if len(genres) > 0 {
		genrePtr = &genres
	}

	var episodes []entity.Episode
	for _, ep := range a.Episodes {
		if ep == nil {
			continue
		}
		episodes = append(episodes, s.mapProtoToEpisode(ep))
	}
	var episodesPtr *[]entity.Episode
	if len(episodes) > 0 {
		episodesPtr = &episodes
	}

	return entity.Anime{
		ID:              uuid.New(),
		Title:           setStr(a.Title),
		TitleJapanese:   setStr(a.TitleJapanese),
		OriginalPageURL: setStr(a.OriginalPageUrl),
		Thumbnail:       setStr(a.Thumbnail),
		Score:           setStr(a.Score),
		Producer:        setStr(a.Producer),
		Type:            setStr(a.Type),
		Status:          setStr(a.Status),
		TotalEpisodes:   setStr(a.TotalEpisodes),
		Duration:        setStr(a.Duration),
		ReleaseDate:     setStr(a.ReleaseDate),
		ReleasedDay:     setStr(a.ReleasedDay),
		Studio:          setStr(a.Studio),
		Genre:           genrePtr,
		Rating:          setStr(a.Rating),
		Episodes:        episodesPtr,
	}
}

func (s *AnimeService) mapProtoToEpisode(ep *pb.AnimeEpisode) entity.Episode {
	setStr := func(v string) *string {
		if v == "" {
			return nil
		}
		vv := v
		return &vv
	}

	var list []entity.ListOfEpisode
	for _, item := range ep.ListEpisode {
		if item == nil {
			continue
		}
		name := setStr(item.Name)
		url := setStr(item.PageUrl)
		list = append(list, entity.ListOfEpisode{Name: name, PageUrl: url})
	}
	var listPtr *[]entity.ListOfEpisode
	if len(list) > 0 {
		listPtr = &list
	}

	var downloads []entity.DownloadLink
	for _, d := range ep.DownloadLinks {
		if d == nil {
			continue
		}
		n := setStr(d.Name)
		u := setStr(d.Url)
		sz := setStr(d.Size)
		q := setStr(d.Quality)
		f := setStr(d.Format)
		downloads = append(downloads, entity.DownloadLink{Name: n, URL: u, Size: sz, Quality: q, Format: f})
	}
	var downloadsPtr *[]entity.DownloadLink
	if len(downloads) > 0 {
		downloadsPtr = &downloads
	}

	return entity.Episode{
		ID:                 uuid.New(),
		Title:              setStr(ep.Title),
		PlayerURL:          setStr(ep.PlayerUrl),
		PageURL:            setStr(ep.PageUrl),
		PostedBy:           setStr(ep.PostedBy),
		PreviousEpisodeURL: setStr(ep.PreviousEpisodeUrl),
		NextEpisodeURL:     setStr(ep.NextEpisodeUrl),
		SeeAllEpisodesURL:  setStr(ep.SeeAllEpisodesUrl),
		ReleaseDate:        setStr(ep.ReleaseDate),
		ReleaseTime:        setStr(ep.ReleaseTime),
		EpisodeNumber:      setStr(ep.EpisodeNumber),
		ListEpisode:        listPtr,
		DownloadLinks:      downloadsPtr,
	}
}

