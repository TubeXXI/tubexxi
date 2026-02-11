package entity

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID              uuid.UUID  `json:"id" validate:"required,uuid"`
	Title           string     `json:"title" validate:"required"`
	OriginalTitle   *string    `json:"original_title" validate:"omitempty"`
	Thumbnail       *string    `json:"thumbnail" validate:"omitempty"`
	Synopsis        *string    `json:"synopsis" validate:"omitempty"`
	Rating          *float64   `json:"rating" validate:"omitempty"`
	Duration        *int64     `json:"duration" validate:"omitempty"`
	Year            *int32     `json:"year" validate:"omitempty"`
	DatePublished   *time.Time `json:"date_published" validate:"omitempty"`
	LabelQuality    *string    `json:"label_quality" validate:"omitempty"`
	Genre           *string    `json:"genre" validate:"omitempty"`
	OriginalPageUrl *string    `json:"original_page_url" validate:"omitempty"`
}

type MovieDetail struct {
	Movie
	Votes        *int64          `json:"votes" validate:"omitempty"`
	ReleaseDate  *time.Time      `json:"release_date" validate:"omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at" validate:"omitempty"`
	PlayerUrl    *[]PlayerUrl    `json:"player_url" validate:"omitempty"`
	TrailerUrl   *string         `json:"trailer_url" validate:"omitempty"`
	Director     *[]MoviePerson  `json:"director" validate:"omitempty"`
	MovieStar    *[]MoviePerson  `json:"movie_star" validate:"omitempty"`
	Countries    *[]CountryMovie `json:"countries" validate:"omitempty"`
	Genres       *[]Genre        `json:"genres" validate:"omitempty"`
	SimilarMovies *[]Movie        `json:"similar_movies" validate:"omitempty"`
}


type MoviePerson struct {
	Name    *string `json:"name" validate:"omitempty"`
	PageUrl *string `json:"page_url" validate:"omitempty"`
}

type PlayerUrl struct {
	URL  *string `json:"url" validate:"omitempty"`
	Type *string `json:"type" validate:"omitempty"`
}

type CountryMovie struct {
	Name    *string `json:"name" validate:"omitempty"`
	PageUrl *string `json:"page_url" validate:"omitempty"`
}

type Genre struct {
	Name    *string `json:"name" validate:"omitempty"`
	PageUrl *string `json:"page_url" validate:"omitempty"`
}

type HomeScrapperListResponse struct {
	CarouselMovies       []Movie `json:"carousel_movies"`
	NewMovies            []Movie `json:"new_movies"`
	FeaturedSeries       []Movie `json:"featured_series"`
	SeriesUpdates        []Movie `json:"series_updates"`
	TopOfTheMonth        []Movie `json:"top_of_the_month"`
	RecommendationForYou []Movie `json:"recommendation_for_you"`
	WatchWithFamily      []Movie `json:"watch_with_family"`
	LatestActionMovies   []Movie `json:"latest_action_movies"`
	KoreanDramaMarathon  []Movie `json:"korean_drama_marathon"`
	LatestHorrorMovies   []Movie `json:"latest_horror_movies"`
	LatestRomanceMovies  []Movie `json:"latest_romance_movies"`
	LatestComedyMovies   []Movie `json:"latest_comedy_movies"`
	LatestKoreanMovies   []Movie `json:"latest_korean_movies"`
	LatestThailandMovies []Movie `json:"latest_thailand_movies"`
	LatestIndianMovies   []Movie `json:"latest_indian_movies"`
	AllLatestMovies      []Movie `json:"all_latest_movies"`
}

type HomeScrapperResponse struct {
	Key        string  `json:"key"`
	Value      []Movie `json:"value"`
	ViewAllUrl *string `json:"view_all_url" validate:"omitempty"`
}

type Pagination struct {
	CurrentPage int32   `json:"current_page"`
	TotalPage   int32   `json:"total_page"`
	HasNext     bool    `json:"has_next"`
	HasPrev     bool    `json:"has_prev"`
	NextPageUrl *string `json:"next_page_url"`
	PrevPageUrl *string `json:"prev_page_url"`
}

type MovieListResponse struct {
	Movies     []Movie    `json:"movies"`
	Pagination Pagination `json:"pagination"`
}
