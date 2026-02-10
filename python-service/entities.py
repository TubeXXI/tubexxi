from dataclasses import dataclass
from typing import List, Optional
import uuid

@dataclass
class Movie:
    id: str  # UUID string
    title: string
    original_title: Optional[str] = None
    thumbnail: Optional[str] = None
    synopsis: Optional[str] = None
    rating: Optional[float] = None
    duration: Optional[int] = None  # in minutes or seconds? Go struct says int64. HTML has "02:11" etc. I'll store minutes or seconds. Go struct doesn't specify unit. Let's assume seconds or just parse it. The HTML has 'content="PT2M11S"'. I'll try to parse to seconds.
    year: Optional[int] = None
    date_published: Optional[str] = None # ISO string or similar
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
