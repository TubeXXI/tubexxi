from bs4 import BeautifulSoup
import uuid
import re
from typing import List, Optional
from entities import Movie, HomeScrapperResponse, MovieListResponse, Pagination, MovieDetail, PlayerUrl, MoviePerson, CountryMovie, Genre

class BaseScraper:
    def __init__(self, html_content: str, base_url: str = "https://tv8.lk21official.cc"):
        self.soup = BeautifulSoup(html_content, 'html.parser')
        self.base_url = base_url.rstrip('/')

    def _make_absolute_url(self, url: Optional[str]) -> Optional[str]:
        if not url:
            return None
        if url.startswith('http'):
            return url
        if url.startswith('/'):
            return f"{self.base_url}{url}"
        return f"{self.base_url}/{url}"

    def _parse_article(self, article) -> Optional[Movie]:
        try:
            # Title
            title_tag = article.find('h3', class_='poster-title')
            if not title_tag:
                 # Try finding title in .video-title (for similar movies)
                 title_tag = article.find('span', class_='video-title')
            
            title = title_tag.get_text(strip=True) if title_tag else "Unknown"

            # URL
            a_tag = article.find('a', itemprop='url')
            if not a_tag:
                a_tag = article.find('a', href=True)
            
            original_page_url = a_tag['href'] if a_tag else None
            original_page_url = self._make_absolute_url(original_page_url)
            
            # Thumbnail
            img_tag = article.find('img', itemprop='image')
            thumbnail = img_tag['src'] if img_tag else None
            
            # Year
            year_tag = article.find('span', class_='year')
            if not year_tag:
                year_tag = article.find('span', class_='video-year')
            
            year_str = year_tag.get_text(strip=True) if year_tag else None
            year = int(year_str) if year_str and year_str.isdigit() else None

            # Rating
            rating_tag = article.find('span', itemprop='ratingValue')
            rating_str = rating_tag.get_text(strip=True) if rating_tag else None
            rating = float(rating_str) if rating_str else None

            # Duration
            duration_tag = article.find('span', class_='duration')
            duration = None
            if duration_tag:
                dur_text = duration_tag.get_text(strip=True)
                if ':' in dur_text:
                    parts = dur_text.split(':')
                    if len(parts) == 2:
                        duration = int(parts[0]) * 60 + int(parts[1])
                    elif len(parts) == 3:
                        duration = int(parts[0]) * 3600 + int(parts[1]) * 60 + int(parts[2])
            
            # Quality
            quality_tag = article.find('span', class_='label')
            label_quality = quality_tag.get_text(strip=True) if quality_tag else None

            # Genre
            genre_meta = article.find('meta', itemprop='genre')
            genre = genre_meta['content'] if genre_meta else None

            # ID
            movie_id = str(uuid.uuid4())

            return Movie(
                id=movie_id,
                title=title,
                original_title=None,
                thumbnail=thumbnail,
                synopsis=None,
                rating=rating,
                duration=duration,
                year=year,
                date_published=None,
                label_quality=label_quality,
                genre=genre,
                original_page_url=original_page_url
            )
        except Exception:
            return None

class HomeScraper(BaseScraper):
    def scrape(self) -> List[HomeScrapperResponse]:
        results = []
        slider_mappings = {
            "New Movies": "Film Terbaru",
            "Featured Series": "SERIES UNGGULAN",
            "Series Updates": "SERIES UPDATE",
            "Top Of The Month": "TOP BULAN INI",
            "Recommendation For You": "Rekomendasi Untukmu",
            "Watch With Family": "Nonton Bareng Keluarga",
            "Latest Action Movies": "Action Terbaru",
            "Korean Drama Marathon": "Maraton Drakor",
            "Latest Horror Movies": "Horror Terbaru",
            "Latest Romance Movies": "Romance Terbaru",
            "Latest Comedy Movies": "Comedy Terbaru",
            "Latest Korean Movies": "Korea Terbaru",
            "Latest Thailand Movies": "Thailand Terbaru",
            "Latest Indian Movies": "India Terbaru",
        }

        # 1. Process Sliders
        for key, aria_label in slider_mappings.items():
            movies, view_all_url = self._scrape_slider(aria_label)
            if movies:
                results.append(HomeScrapperResponse(key=key, value=movies, view_all_url=view_all_url))
            else:
                results.append(HomeScrapperResponse(key=key, value=[], view_all_url=None))

        # 2. Process "All Latest Movies" (Grid)
        all_latest = self._scrape_grid("Daftar Lengkap Film Terbaru")
        results.append(HomeScrapperResponse(key="All Latest Movies", value=all_latest, view_all_url=self._make_absolute_url("/latest-movies")))

        return results

    def _scrape_slider(self, aria_label: str) -> tuple[List[Movie], Optional[str]]:
        slider_wrapper = self.soup.find('div', class_='slider-wrapper', attrs={'aria-label': aria_label})
        if not slider_wrapper:
            return [], None
        
        view_all_url = None
        section = slider_wrapper.find_parent('section')
        if section:
            see_all_link = section.find('a', string=re.compile(r"LIHAT SEMUA", re.IGNORECASE))
            if see_all_link:
                view_all_url = self._make_absolute_url(see_all_link['href'])

        movies = []
        articles = slider_wrapper.find_all('article')
        for article in articles:
            movie = self._parse_article(article)
            if movie:
                movies.append(movie)
        return movies, view_all_url

    def _scrape_grid(self, header_text: str) -> List[Movie]:
        header = self.soup.find('h2', string=header_text)
        if not header:
            return []
        
        header_div = header.find_parent('div', class_='header')
        if not header_div:
            return []
        
        gallery_grid = header_div.find_next_sibling('div', class_='gallery-grid')
        if not gallery_grid:
            return []

        movies = []
        articles = gallery_grid.find_all('article')
        for article in articles:
            movie = self._parse_article(article)
            if movie:
                movies.append(movie)
        return movies

class MovieListScraper(BaseScraper):
    def scrape(self) -> MovieListResponse:
        movies = []
        
        gallery_grid = self.soup.find('div', class_='gallery-grid')
        if not gallery_grid:
             return MovieListResponse(movies=[], pagination=Pagination(current_page=1, total_page=1, has_next=False, has_prev=False))
        
        articles = gallery_grid.find_all('article')
        for article in articles:
            movie = self._parse_article(article)
            if movie:
                movies.append(movie)
                
        # Pagination Logic
        pagination_div = self.soup.find('div', class_='pagination')
        current_page = 1
        total_page = 1
        has_next = False
        has_prev = False
        next_page_url = None
        prev_page_url = None

        if pagination_div:
            # Current Page
            current_span = pagination_div.find('span', class_='current')
            if current_span:
                try:
                    current_page = int(current_span.get_text(strip=True))
                except:
                    pass
            
            # Next Page
            next_link = pagination_div.find('a', class_='next page-numbers')
            if next_link:
                has_next = True
                next_page_url = self._make_absolute_url(next_link['href'])
            
            # Prev Page
            prev_link = pagination_div.find('a', class_='prev page-numbers')
            if prev_link:
                has_prev = True
                prev_page_url = self._make_absolute_url(prev_link['href'])

            # Total Page
            page_numbers = pagination_div.find_all('a', class_='page-numbers')
            numeric_pages = [p for p in page_numbers if p.get_text(strip=True).isdigit()]
            if numeric_pages:
                last_page_num = int(numeric_pages[-1].get_text(strip=True))
                if last_page_num > total_page:
                    total_page = last_page_num
            
            if current_page > total_page:
                total_page = current_page

        return MovieListResponse(
            movies=movies,
            pagination=Pagination(
                current_page=current_page,
                total_page=total_page,
                has_next=has_next,
                has_prev=has_prev,
                next_page_url=next_page_url,
                prev_page_url=prev_page_url
            )
        )

class MovieDetailScraper(BaseScraper):
    def scrape(self, original_url: str) -> Optional[MovieDetail]:
        # Basic Info
        title_h1 = self.soup.find('h1', class_='entry-title')
        title = title_h1.get_text(strip=True) if title_h1 else "Unknown"
        
        # Synopsis
        synopsis = ""
        entry_content = self.soup.find('div', class_='entry-content')
        if entry_content:
            synopsis = entry_content.get_text(strip=True)
        
        # Thumbnail
        thumbnail = None
        img_tag = self.soup.find('img', class_='attachment-post-thumbnail')
        if img_tag:
            thumbnail = img_tag.get('src')
        
        # Meta Data (Rating, Year, Duration, etc.)
        # Usually found in .gmr-movie-data or similar
        
        rating = None
        rating_tag = self.soup.find('span', class_='gmr-rating-value') # adjust selector if needed
        if not rating_tag:
             rating_tag = self.soup.find('span', itemprop='ratingValue')
        if rating_tag:
            try:
                rating = float(rating_tag.get_text(strip=True))
            except:
                pass
        
        votes = None
        votes_tag = self.soup.find('span', itemprop='ratingCount')
        if votes_tag:
            try:
                votes = int(votes_tag.get_text(strip=True).replace(',', ''))
            except:
                pass

        year = None
        year_link = self.soup.find('a', href=re.compile(r'/year/'))
        if year_link:
            try:
                year = int(year_link.get_text(strip=True))
            except:
                pass
        
        duration = None
        duration_tag = self.soup.find('span', itemprop='duration') # content="PT..."
        if duration_tag:
             # Try parsing content attribute or text
             # If content is PT2M11S, we need parsing. 
             # Or look for text like "2h 4m"
             dur_text = duration_tag.get_text(strip=True)
             # Simple parsing logic for "2h 4m" or "124 min"
             total_seconds = 0
             h_match = re.search(r'(\d+)h', dur_text)
             m_match = re.search(r'(\d+)m', dur_text)
             if h_match:
                 total_seconds += int(h_match.group(1)) * 3600
             if m_match:
                 total_seconds += int(m_match.group(1)) * 60
             if total_seconds > 0:
                 duration = total_seconds
        
        # Player URLs
        player_urls = []
        player_list = self.soup.find('ul', id='player-list')
        if player_list:
            for li in player_list.find_all('li'):
                a = li.find('a')
                if a:
                    url = a.get('data-url') or a.get('href')
                    server = a.get('data-server') or a.get_text(strip=True)
                    if url:
                        player_urls.append(PlayerUrl(url=url, type=server))
        
        # Directors
        directors = []
        # Find section for director
        # This is tricky without exact HTML for meta section. Assuming links with rel="tag" or specific structure
        # Often: Director: <a href="...">Name</a>
        
        # Stars
        movie_stars = []
        
        # Countries
        countries = []
        
        # Genres
        genres = []
        
        # General Meta Parsing (loop through gmr-movie-inn or similar)
        # Based on example site structure (not provided in full, but typical WP):
        meta_items = self.soup.find_all('div', class_='gmr-moviedata')
        for item in meta_items:
            text = item.get_text(strip=True)
            if "Director" in text or "Sutradara" in text:
                for a in item.find_all('a'):
                    directors.append(MoviePerson(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
            elif "Cast" in text or "Bintang" in text or "Pemeran" in text:
                for a in item.find_all('a'):
                    movie_stars.append(MoviePerson(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
            elif "Country" in text or "Negara" in text:
                for a in item.find_all('a'):
                    countries.append(CountryMovie(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
            elif "Genre" in text:
                for a in item.find_all('a'):
                    genres.append(Genre(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
            elif "Release" in text or "Rilis" in text:
                 # Extract date
                 pass

        # Similar Movies
        similar_movies = []
        related_content = self.soup.find('div', class_='related-content')
        if related_content:
            for article in related_content.find_all('li'): # It's ul > li based on example
                 # Need to adapt _parse_article or write new logic because structure might differ (li vs article)
                 # Example: <li> <a href...><picture>... <div class="video-info">...
                 a_tag = article.find('a')
                 if a_tag:
                     sim_url = self._make_absolute_url(a_tag['href'])
                     sim_title_tag = article.find('span', class_='video-title')
                     sim_title = sim_title_tag.get_text(strip=True) if sim_title_tag else "Unknown"
                     sim_thumb_tag = article.find('img')
                     sim_thumb = sim_thumb_tag['src'] if sim_thumb_tag else None
                     sim_year_tag = article.find('span', class_='video-year')
                     sim_year = int(sim_year_tag.get_text(strip=True)) if sim_year_tag and sim_year_tag.get_text(strip=True).isdigit() else None
                     
                     similar_movies.append(Movie(
                         id=str(uuid.uuid4()),
                         title=sim_title,
                         original_page_url=sim_url,
                         thumbnail=sim_thumb,
                         year=sim_year
                     ))

        # ID
        movie_id = str(uuid.uuid4())

        movie = Movie(
            id=movie_id,
            title=title,
            thumbnail=thumbnail,
            synopsis=synopsis,
            rating=rating,
            duration=duration,
            year=year,
            original_page_url=original_url,
            # Populate genre string from list
            genre=", ".join([g.name for g in genres]) if genres else None
        )

        return MovieDetail(
            movie=movie,
            votes=votes,
            player_urls=player_urls,
            directors=directors,
            movie_stars=movie_stars,
            countries=countries,
            genres=genres,
            similar_movies=similar_movies
        )
