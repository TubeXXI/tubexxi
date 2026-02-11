from bs4 import BeautifulSoup
import uuid
import re
from typing import List, Optional
from datetime import datetime
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
        title = "Unknown"
        title_h1 = self.soup.find('div', class_='movie-info')
        if title_h1:
            h1 = title_h1.find('h1')
            if h1:
                # Clean title: "Nonton The Wrecking Crew (2026) Sub Indo di Lk21" -> "The Wrecking Crew (2026)"
                raw_title = h1.get_text(strip=True)
                title = raw_title.replace("Nonton ", "").replace(" Sub Indo di Lk21", "")
        
        # Synopsis
        synopsis = ""
        synopsis_div = self.soup.find('div', class_='synopsis')
        if synopsis_div:
            synopsis = synopsis_div.get('data-full') or synopsis_div.get_text(strip=True)
        
        # Meta Data from info-tag (Rating, Quality, Duration)
        rating = None
        label_quality = None
        duration = None
        
        info_tag = self.soup.find('div', class_='info-tag')
        if info_tag:
            spans = info_tag.find_all('span', recursive=False)
            # 1. Rating
            if len(spans) > 0:
                try:
                    rating_text = spans[0].get_text(strip=True)
                    rating = float(rating_text)
                except:
                    pass
            
            # 2. Quality (e.g., WEBDL)
            if len(spans) > 1:
                label_quality = spans[1].get_text(strip=True)
            
            # 3. Resolution (e.g., 1080p) - optional, maybe append to quality?
            
            # 4. Duration (e.g., 2h 4m)
            if len(spans) > 3:
                dur_text = spans[3].get_text(strip=True)
                total_seconds = 0
                h_match = re.search(r'(\d+)h', dur_text)
                m_match = re.search(r'(\d+)m', dur_text)
                if h_match:
                    total_seconds += int(h_match.group(1)) * 3600
                if m_match:
                    total_seconds += int(m_match.group(1)) * 60
                if total_seconds > 0:
                    duration = total_seconds

        # Genres & Countries from tag-list
        genres = []
        countries = []
        tag_list = self.soup.find('div', class_='tag-list')
        if tag_list:
            for a in tag_list.find_all('a'):
                href = a.get('href', '')
                text = a.get_text(strip=True)
                if '/genre/' in href:
                    genres.append(Genre(name=text, page_url=self._make_absolute_url(href)))
                elif '/country/' in href:
                    countries.append(CountryMovie(name=text, page_url=self._make_absolute_url(href)))

        # Detailed Metadata (Directors, Cast, Votes, Release, Updated, Thumbnail)
        directors = []
        movie_stars = []
        votes = None
        release_date = None
        updated_at = None
        thumbnail = None
        
        detail_div = self.soup.find('div', class_='detail')
        if detail_div:
            # Thumbnail inside detail div
            img_tag = detail_div.find('img', itemprop='image')
            if img_tag:
                thumbnail = img_tag.get('src')
            
            # Parse paragraphs
            for p in detail_div.find_all('p'):
                text = p.get_text(strip=True)
                
                if "Sutradara:" in text:
                    for a in p.find_all('a'):
                        directors.append(MoviePerson(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
                
                elif "Bintang Film:" in text:
                    for a in p.find_all('a'):
                        movie_stars.append(MoviePerson(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
                
                elif "Negara:" in text and not countries: # Fallback if not found in tag-list
                    for a in p.find_all('a'):
                        countries.append(CountryMovie(name=a.get_text(strip=True), page_url=self._make_absolute_url(a['href'])))
                
                elif "Votes:" in text:
                    try:
                        # "Votes: 170"
                        votes_str = text.replace("Votes:", "").strip().replace(',', '')
                        votes = int(votes_str)
                    except:
                        pass
                
                elif "Release:" in text:
                    release_date = text.replace("Release:", "").strip()
                
                elif "Updated:" in text:
                    updated_at = text.replace("Updated:", "").strip()

        # Year (extract from title or release date if missing)
        year = None
        # Try from Release Date
        if release_date:
             year_match = re.search(r'\d{4}', release_date)
             if year_match:
                 year = int(year_match.group(0))
        # Fallback to title
        if not year:
            year_match = re.search(r'\((\d{4})\)', title)
            if year_match:
                year = int(year_match.group(1))

        # Date Parsing to ISO
        iso_release_date = None
        if release_date:
             # "28 Jan 2026"
             try:
                 dt = datetime.strptime(release_date, "%d %b %Y")
                 iso_release_date = dt.strftime("%Y-%m-%dT%H:%M:%SZ")
             except:
                 pass
        
        iso_updated_at = None
        if updated_at:
             # "01 Feb 2026 14:58:14"
             try:
                 dt = datetime.strptime(updated_at, "%d %b %Y %H:%M:%S")
                 iso_updated_at = dt.strftime("%Y-%m-%dT%H:%M:%SZ")
             except:
                 pass

        # Trailer URL
        trailer_url = None
        trailer_link = self.soup.find('a', class_='yt-lightbox')
        if trailer_link:
            trailer_url = trailer_link.get('href')

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
        
        # Similar Movies
        similar_movies = []
        related_content = self.soup.find('div', class_='related-content')
        if related_content:
             video_list = related_content.find('ul', class_='video-list')
             if video_list:
                for li in video_list.find_all('li'):
                     a_tag = li.find('a')
                     if a_tag:
                         sim_url = self._make_absolute_url(a_tag['href'])
                         
                         sim_title = "Unknown"
                         sim_title_tag = li.find('span', class_='video-title')
                         if sim_title_tag:
                             sim_title = sim_title_tag.get_text(strip=True)
                         
                         sim_thumb = None
                         sim_thumb_tag = li.find('img')
                         if sim_thumb_tag:
                             sim_thumb = sim_thumb_tag.get('src')
                         
                         sim_year = None
                         sim_year_tag = li.find('span', class_='video-year')
                         if sim_year_tag:
                             try:
                                 sim_year = int(sim_year_tag.get_text(strip=True))
                             except:
                                 pass
                         
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
            label_quality=label_quality,
            date_published=iso_release_date, # Use release date as date_published
            genre=", ".join([g.name for g in genres]) if genres else None
        )

        return MovieDetail(
            movie=movie,
            votes=votes,
            release_date=iso_release_date,
            updated_at=iso_updated_at,
            player_urls=player_urls,
            trailer_url=trailer_url,
            directors=directors,
            movie_stars=movie_stars,
            countries=countries,
            genres=genres,
            similar_movies=similar_movies
        )

