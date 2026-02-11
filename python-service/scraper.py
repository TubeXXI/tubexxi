from bs4 import BeautifulSoup
import json
import uuid
import re
from typing import List, Optional
from datetime import datetime
from entities import Movie, HomeScrapperResponse, MovieListResponse, Pagination, MovieDetail, PlayerUrl, MoviePerson, CountryMovie, Genre, SeriesDetail, SeasonList, EpisodeList, SeriesEpisodeResponse

class BaseScraper:
    def __init__(self, html_content: str, base_url: str = "https://tv8.lk21official.cc"):
        self.soup = BeautifulSoup(html_content, 'html.parser')
        self.base_url = base_url.rstrip('/')

    def _make_absolute_url(self, url: Optional[str]) -> Optional[str]:
        if url is None:
            return None
        
        # Ensure it is a string before stripping
        try:
             url = str(url).strip()
        except Exception:
             return None

        if not url:
            return None

        if url.startswith('http'):
            return url
        
        if url.startswith('//'):
            return f"https:{url}"
        
        # Ensure base_url is set
        base = self.base_url if self.base_url else "https://tv8.lk21official.cc"
        base = base.rstrip('/') # Ensure no trailing slash
        
        # Safe concatenation
        if not url.startswith('/'):
            url = f"/{url}"
            
        return f"{base}{url}"

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
            "New Movies": "TERBARU",
            "Featured Series": "Film Unggulan",
            "Series Updates": "LK21 TERBARU",
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
            
            # If empty, try fallback to grid headers for Series Home
            if not movies and self.base_url and "nontondrama" in self.base_url:
                 if key == "New Movies":
                      # For "New Movies", fallback to "Episode Terbaru" grid
                      movies = self._scrape_grid("Episode Terbaru")
                 elif key == "Featured Series":
                      pass
            
            # Safe check for view_all_url
            final_view_all_url = None
            try:
                final_view_all_url = self._make_absolute_url(view_all_url)
            except:
                pass
            
            # Ensure safe URL for movies too
            safe_movies = []
            for m in movies:
                  try:
                      m.original_page_url = self._make_absolute_url(m.original_page_url)
                      m.thumbnail = self._make_absolute_url(m.thumbnail)
                      safe_movies.append(m)
                  except:
                       pass

            if safe_movies:
                results.append(HomeScrapperResponse(key=key, value=safe_movies, view_all_url=final_view_all_url))
            else:
                results.append(HomeScrapperResponse(key=key, value=[], view_all_url=final_view_all_url))

        # 2. Process "All Latest Movies" (Grid)
        all_latest = self._scrape_grid("Daftar Lengkap Series Terbaru")
        if not all_latest:
             all_latest = self._scrape_grid("Daftar Lengkap Film Terbaru")
             
        results.append(HomeScrapperResponse(key="All Latest Movies", value=all_latest, view_all_url=self._make_absolute_url("/latest-movies")))

        return results

    def _scrape_slider(self, aria_label: str) -> tuple[List[Movie], Optional[str]]:
        slider_wrapper = self.soup.find('div', class_='slider-wrapper', attrs={'aria-label': aria_label})
        
        # Fallback: Check if aria_label is actually in a section header h2/h3
        if not slider_wrapper:
             # Look for section with header containing aria_label
             headers = self.soup.find_all(['h2', 'h3'], string=re.compile(re.escape(aria_label), re.IGNORECASE))
             for header in headers:
                 # Find parent section or container
                 section = header.find_parent('section')
                 if section:
                     slider_wrapper = section.find('div', class_='slider-wrapper')
                     if slider_wrapper:
                         break
                     # Or maybe it's a grid
                     grid = section.find('div', class_='gallery-grid')
                     if grid:
                         # Treat grid as slider wrapper for parsing purposes (both contain articles)
                         slider_wrapper = grid
                         break

        if not slider_wrapper:
            return [], None
        
        view_all_url = None
        # Try finding view all link in parent section
        section = slider_wrapper.find_parent('section')
        if section:
            see_all_link = section.find('a', string=re.compile(r"SEMUA", re.IGNORECASE))
            if see_all_link:
                view_all_url = self._make_absolute_url(see_all_link['href'])
        
        # If not found in section, maybe it's a direct sibling or inside header div
        if not view_all_url:
             header_div = slider_wrapper.find_previous_sibling('div', class_='header')
             if header_div:
                  see_all_link = header_div.find('a', string=re.compile(r"SEMUA", re.IGNORECASE))
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
                
        current_page = 1
        total_page = 1
        has_next = False
        has_prev = False
        next_page_url = None
        prev_page_url = None

        def extract_page_number(href: Optional[str]) -> Optional[int]:
            if not href:
                return None
            m = re.search(r"/page/(\d+)", href)
            if not m:
                return None
            try:
                return int(m.group(1))
            except Exception:
                return None

        pagination_div = self.soup.find('div', class_='pagination')
        if pagination_div:
            current_span = pagination_div.find('span', class_='current')
            if current_span:
                try:
                    current_page = int(current_span.get_text(strip=True))
                except Exception:
                    pass

            next_link = pagination_div.find('a', class_='next page-numbers')
            if next_link and next_link.get('href'):
                has_next = True
                next_page_url = self._make_absolute_url(next_link['href'])

            prev_link = pagination_div.find('a', class_='prev page-numbers')
            if prev_link and prev_link.get('href'):
                has_prev = True
                prev_page_url = self._make_absolute_url(prev_link['href'])

            page_numbers = pagination_div.find_all('a', class_='page-numbers')
            max_page = 1
            for p in page_numbers:
                txt = p.get_text(strip=True)
                href = p.get('href')
                if txt.isdigit():
                    try:
                        max_page = max(max_page, int(txt))
                    except Exception:
                        pass
                page_num = extract_page_number(href)
                if page_num:
                    max_page = max(max_page, page_num)
            total_page = max(total_page, max_page, current_page)
        else:
            pagination_ul = None
            wrapper = self.soup.find('nav', class_='pagination-wrapper')
            if wrapper:
                pagination_ul = wrapper.find('ul', class_='pagination')
            if not pagination_ul:
                pagination_ul = self.soup.find('ul', class_='pagination')

            if pagination_ul:
                links = pagination_ul.find_all('a', href=True)
                page_to_url: dict[int, str] = {}
                max_page = 1
                for a in links:
                    href = self._make_absolute_url(a.get('href'))
                    if not href:
                        continue
                    page_num = extract_page_number(href)
                    if page_num:
                        page_to_url[page_num] = href
                        max_page = max(max_page, page_num)
                    txt = a.get_text(strip=True)
                    if txt.isdigit():
                        try:
                            max_page = max(max_page, int(txt))
                        except Exception:
                            pass
                total_page = max(total_page, max_page)

                active_li = pagination_ul.find('li', class_=re.compile(r"\bactive\b", re.IGNORECASE))
                if active_li:
                    active_a = active_li.find('a')
                    if active_a:
                        txt = active_a.get_text(strip=True)
                        if txt.isdigit():
                            current_page = int(txt)
                        else:
                            page_num = extract_page_number(active_a.get('href'))
                            if page_num:
                                current_page = page_num
                current_page = max(1, current_page)
                total_page = max(total_page, current_page)

                if (current_page + 1) in page_to_url:
                    has_next = True
                    next_page_url = page_to_url[current_page+1]
                elif current_page < total_page:
                    has_next = True
                if (current_page - 1) in page_to_url:
                    has_prev = True
                    prev_page_url = page_to_url[current_page-1]
                elif current_page > 1:
                    has_prev = True

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

class SeriesDetailScraper(MovieDetailScraper):
    def scrape(self, original_url: str) -> Optional[SeriesDetail]:
        # 1. Get Base Movie Detail
        movie_detail = super().scrape(original_url)
        if not movie_detail:
            return None
        
        season_list: List[SeasonList] = []
        season_name: Optional[str] = None
        status: Optional[str] = None

        status_tag = self.soup.find(string=re.compile(r"Status\s*:\s*", re.IGNORECASE))
        if status_tag:
            parent = getattr(status_tag, "parent", None)
            if parent:
                status = parent.get_text(strip=True)
                status = re.sub(r"Status\s*:\s*", "", status, flags=re.IGNORECASE).strip() or None

        episodes_by_season: dict[int, dict[int, str]] = {}
        current_season_num: Optional[int] = None

        season_data_el = self.soup.find(id="season-data")
        if season_data_el:
            raw = season_data_el.get_text(strip=True)
            if raw:
                try:
                    season_data = json.loads(raw)
                    if isinstance(season_data, dict):
                        for season_key, eps in season_data.items():
                            try:
                                season_num = int(str(season_key))
                            except Exception:
                                continue
                            if not isinstance(eps, list):
                                continue
                            for ep in eps:
                                if not isinstance(ep, dict):
                                    continue
                                try:
                                    episode_num = int(str(ep.get("episode_no") or ep.get("episode") or "0"))
                                except Exception:
                                    continue
                                slug = ep.get("slug")
                                if not slug:
                                    continue
                                url = self._make_absolute_url("/" + str(slug).lstrip("/"))
                                episodes_by_season.setdefault(season_num, {})[episode_num] = url
                except Exception:
                    pass

        watch_data_el = self.soup.find(id="watch-history-data")
        if watch_data_el:
            raw = watch_data_el.get_text(strip=True)
            if raw:
                try:
                    watch_data = json.loads(raw)
                    if isinstance(watch_data, dict) and watch_data.get("current_season") is not None:
                        current_season_num = int(str(watch_data.get("current_season")))
                except Exception:
                    pass

        if not episodes_by_season:
            for a in self.soup.find_all("a", href=True):
                href = a.get("href")
                if not href:
                    continue
                m = re.search(r"season-(\d+)-episode-(\d+)", href, re.IGNORECASE)
                if not m:
                    continue
                try:
                    season_num = int(m.group(1))
                    episode_num = int(m.group(2))
                except Exception:
                    continue
                episodes_by_season.setdefault(season_num, {})[episode_num] = self._make_absolute_url(href)

        if not episodes_by_season:
            current_season_num = 1
            total_seasons = 1

            season_select = self.soup.find("select", class_="season-select")
            if season_select:
                options = season_select.find_all("option")
                if options:
                    total_seasons = len(options)
                selected_option = season_select.find("option", selected=True)
                if selected_option:
                    try:
                        current_season_num = int(selected_option.get("value") or "1")
                        season_name = selected_option.get_text(strip=True) or None
                    except Exception:
                        pass

            episodes: List[EpisodeList] = []
            episode_ul = self.soup.find("ul", class_="episode-list")
            if episode_ul:
                for li in episode_ul.find_all("li"):
                    a = li.find("a", href=True)
                    if not a:
                        continue
                    ep_url = self._make_absolute_url(a.get("href"))
                    ep_title = a.get_text(strip=True)
                    ep_num = 0
                    match = re.search(r"Episode\s+(\d+)", ep_title or "", re.IGNORECASE)
                    if match:
                        ep_num = int(match.group(1))
                    episodes.append(EpisodeList(
                        episode_number=ep_num,
                        episode_url=ep_url,
                        player_urls=[],
                        trailer_url=None,
                    ))

            season_list.append(SeasonList(
                current_season=current_season_num,
                total_season=total_seasons,
                episode_list=episodes,
            ))
        else:
            seasons_sorted = sorted(episodes_by_season.keys())
            total_seasons = max(seasons_sorted)
            if current_season_num is None:
                current_season_num = total_seasons
            season_name = f"Season {current_season_num}"

            for season_num in seasons_sorted:
                episode_map = episodes_by_season.get(season_num, {})
                episodes: List[EpisodeList] = []
                for episode_num in sorted(episode_map.keys()):
                    episodes.append(EpisodeList(
                        episode_number=episode_num,
                        episode_url=episode_map[episode_num],
                        player_urls=[],
                        trailer_url=None,
                    ))
                season_list.append(SeasonList(
                    current_season=season_num,
                    total_season=total_seasons,
                    episode_list=episodes,
                ))
        
        return SeriesDetail(
            movie_detail=movie_detail,
            season_name=season_name,
            status=status,
            season_list=season_list
        )

class SeriesEpisodeScraper(BaseScraper):
    def scrape(self, url: str) -> Optional[SeriesEpisodeResponse]:
        # Parse Episode Number from Title or URL
        title_tag = self.soup.find('title')
        episode_number = 0
        if title_tag:
            title_text = title_tag.get_text(strip=True)
            match = re.search(r'Episode\s+(\d+)', title_text, re.IGNORECASE)
            if match:
                episode_number = int(match.group(1))
        
        # Player URLs
        player_urls = []
        player_list = self.soup.find('ul', id='player-list')
        if player_list:
            for li in player_list.find_all('li'):
                a = li.find('a')
                if a:
                    p_url = a.get('data-url') or a.get('href')
                    server = a.get('data-server') or a.get_text(strip=True)
                    if p_url:
                        player_urls.append(PlayerUrl(url=p_url, type=server))
        
        # Trailer
        trailer_url = None
        trailer_link = self.soup.find('a', class_='yt-lightbox')
        if trailer_link:
            trailer_url = trailer_link.get('href')

        # Download URL
        download_url = None
        # Try to find download link
        download_btn = self.soup.find('a', class_='btn-download')
        if download_btn:
             download_url = self._make_absolute_url(download_btn['href'])

        return SeriesEpisodeResponse(
            episode_number=episode_number,
            player_urls=player_urls,
            trailer_url=trailer_url,
            download_url=download_url
        )
