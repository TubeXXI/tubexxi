from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Movie:
    id: str  # UUID string
    title: str
    original_title: Optional[str] = None
    thumbnail: Optional[str] = None
    synopsis: Optional[str] = None
    rating: Optional[float] = None
    duration: Optional[int] = None  # seconds
    year: Optional[int] = None
    date_published: Optional[str] = None # ISO string
    label_quality: Optional[str] = None
    genre: Optional[str] = None
    original_page_url: Optional[str] = None

@dataclass
class Pagination:
    current_page: int
    total_page: int
    has_next: bool
    has_prev: bool
    next_page_url: Optional[str] = None
    prev_page_url: Optional[str] = None

@dataclass
class MovieListResponse:
    movies: List[Movie]
    pagination: Pagination

@dataclass
class HomeScrapperResponse:
    key: str
    value: List[Movie]
    view_all_url: Optional[str] = None

@dataclass
class PlayerUrl:
    url: Optional[str] = None
    type: Optional[str] = None

@dataclass
class MoviePerson:
    name: Optional[str] = None
    page_url: Optional[str] = None

@dataclass
class CountryMovie:
    name: Optional[str] = None
    page_url: Optional[str] = None

@dataclass
class Genre:
    name: Optional[str] = None
    page_url: Optional[str] = None

@dataclass
class MovieDetail:
    movie: Movie
    votes: Optional[int] = None
    release_date: Optional[str] = None
    updated_at: Optional[str] = None
    player_urls: Optional[List[PlayerUrl]] = None
    trailer_url: Optional[str] = None
    directors: Optional[List[MoviePerson]] = None
    movie_stars: Optional[List[MoviePerson]] = None
    countries: Optional[List[CountryMovie]] = None
    genres: Optional[List[Genre]] = None
    similar_movies: Optional[List[Movie]] = None

@dataclass
class SeriesEpisodeResponse:
    episode_number: int
    player_urls: List[PlayerUrl]
    trailer_url: Optional[str] = None
    download_url: Optional[str] = None

@dataclass
class EpisodeList:
    episode_number: int
    episode_url: str
    player_urls: Optional[List[PlayerUrl]] = None
    trailer_url: Optional[str] = None

@dataclass
class SeasonList:
    current_season: int
    total_season: int
    episode_list: List[EpisodeList]

@dataclass
class SeriesDetail:
    movie_detail: MovieDetail
    season_name: Optional[str] = None
    status: Optional[str] = None
    season_list: Optional[List[SeasonList]] = None
