from bs4 import BeautifulSoup
import requests
import json
import dataclasses
import random
import time
import base64
from fake_useragent import UserAgent
from dataclasses import dataclass, field
from typing import List, Optional
import os
import re
from urllib.parse import urlparse, parse_qs, urlencode
try:
    from dotenv import load_dotenv
except ImportError:
    def load_dotenv(*args, **kwargs):
        pass

load_dotenv()

class EnhancedJSONEncoder(json.JSONEncoder):
    def default(self, o):
        if dataclasses.is_dataclass(o):
            return dataclasses.asdict(o)
        return super().default(o)

@dataclass
class Anime:
    title: Optional[str] = None
    title_japanese: Optional[str] = None
    original_page_url: Optional[str] = None
    thumbnail: Optional[str] = None
    score: Optional[str] = None
    producer: Optional[str] = None
    type: Optional[str] = None
    status: Optional[str] = None
    total_episodes: Optional[str] = None
    duration: Optional[str] = None
    release_date: Optional[str] = None
    released_day: Optional[str] = None
    studio: Optional[str] = None
    genre: Optional[List[Genre]] = field(default_factory=list)
    rating: Optional[str] = None
    episodes: Optional[List[Episode]] = field(default_factory=list)

@dataclass
class PaginationInfo:
    total_pages: int
    current_page: int
    has_next: bool
    has_previous: bool
    next_page_url: Optional[str] = None
    previous_page_url: Optional[str] = None
    page_numbers: List[int] = field(default_factory=list)
    per_page: int = 20  # Default items per page

@dataclass
class ScrapeResult:
    animes: List[Anime]
    pagination: PaginationInfo
    query: Optional[str] = None
    
@dataclass
class Genre:
    name: Optional[str] = None
    url: Optional[str] = None
    
@dataclass
class Episode:
    title: Optional[str] = None
    player_url: Optional[str] = None
    page_url: Optional[str] = None
    posted_by: Optional[str] = None
    previous_episode_url: Optional[str] = None
    next_episode_url: Optional[str] = None
    see_all_episodes_url: Optional[str] = None
    release_date: Optional[str] = None
    release_time: Optional[str] = None
    list_episode: Optional[List[ListOfEpisode]] = field(default_factory=list)
    download_links: Optional[List[DownloadLink]] = field(default_factory=list)
    episode_number: Optional[str] = None

@dataclass
class ListOfEpisode:
    name: Optional[str] = None
    page_url: Optional[str] = None

@dataclass
class DownloadLink:
    name: Optional[str] = None
    url: Optional[str] = None
    size: Optional[str] = None
    quality: Optional[str] = None
    format: Optional[str] = None

class OtakuDesuScraper: 
    def __init__(self, base_url="https://otakudesu.best/"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.ua = UserAgent()
        self.timeout = 15
        self.retry_count = 3
    
    def _get_headers(self):
        return {
            'User-Agent': self.ua.random,
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'id,en-US;q=0.7,en;q=0.3',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Referer': self.base_url,
            'Upgrade-Insecure-Requests': '1',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'none',
            'Cache-Control': 'max-age=0',
        }
    
    def _make_request(self, url, max_retries=3):
        """Enhanced request dengan retry, delay, dan rotasi user-agent"""
        for attempt in range(max_retries):
            try:
                # Rotasi User-Agent setiap request
                headers = self._get_headers()
                
                # Random delay antara 2-5 detik
                time.sleep(random.uniform(2, 5))
                
                response = self.session.get(
                    url, 
                    headers=headers, 
                    timeout=self.timeout,
                    allow_redirects=True
                )
                
                if response.status_code == 200:
                    return response
                elif response.status_code == 403:
                    print(f"⚠️ Attempt {attempt + 1}: Blocked (403) for {url}")
                    
                    # Ganti IP dengan proxy (jika punya)
                    # self.session.proxies = {'http': 'http://new-proxy:port'}
                    
                    if attempt < max_retries - 1:
                        time.sleep(random.uniform(5, 10))
                        continue
                
                response.raise_for_status()
                
            except requests.exceptions.RequestException as e:
                print(f"⚠️ Attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(random.uniform(3, 7))
                else:
                    return None
        
        return None
    
    def _parse_pagination(self, soup):
        """Parse pagination info dari BeautifulSoup object"""
        pagination_div = soup.find('div', class_='pagination')
        if not pagination_div:
            return None
        
        # Parse "Pages X of Y"
        total_pages = 1
        current_page = 1
        naviright = pagination_div.find('span', class_='naviright')
        if naviright:
            page_text = naviright.text.strip()
            match = re.search(r'Pages?\s*(\d+)\s*of\s*(\d+)', page_text, re.IGNORECASE)
            if match:
                current_page = int(match.group(1))
                total_pages = int(match.group(2))
        
        # Parse navigasi links
        next_page_url = None
        previous_page_url = None
        page_numbers = []
        
        navileft = pagination_div.find('span', class_='navileft')
        if navileft:
            # Next page
            next_link = navileft.find('a', string='»')
            if next_link and next_link.get('href'):
                next_page_url = next_link.get('href')
            
            # Previous page
            prev_link = navileft.find('a', string='«')
            if prev_link and prev_link.get('href'):
                previous_page_url = prev_link.get('href')
            
            # Page numbers
            for a in navileft.find_all('a'):
                text = a.text.strip()
                if text.isdigit():
                    page_numbers.append(int(text))
        
        # Check for next/previous from rel="next/prev"
        if not next_page_url:
            link_next = soup.find('link', rel='next')
            if link_next and link_next.get('href'):
                next_page_url = link_next.get('href')
        
        if not previous_page_url:
            link_prev = soup.find('link', rel='prev')
            if link_prev and link_prev.get('href'):
                previous_page_url = link_prev.get('href')
        
        return PaginationInfo(
            total_pages=total_pages,
            current_page=current_page,
            has_next=next_page_url is not None,
            has_previous=previous_page_url is not None,
            next_page_url=next_page_url,
            previous_page_url=previous_page_url,
            page_numbers=sorted(page_numbers),
            per_page=len(self._extract_animes(soup)) or 20
        )
        
    def _extract_animes(self, soup):
        """Extract animes dari BeautifulSoup object"""
        animes = []
        
        # Cari ul.chivsrc
        ul_item = soup.find('ul', class_='chivsrc')
        if not ul_item:
            return animes
        
        li_items = ul_item.find_all('li', recursive=False)
        
        for li in li_items:
            # Skip jika ini pagination
            if li.find('div', class_='pagination'):
                continue
            
            # Extract thumbnail - coba ambil dari src atau srcset
            thumbnail = None
            img = li.find('img')
            if img:
                # Prioritaskan srcset untuk kualitas lebih baik
                if img.get('srcset'):
                    srcset_parts = img.get('srcset').split(',')[0].strip().split(' ')
                    if srcset_parts:
                        thumbnail = srcset_parts[0]
                elif img.get('src'):
                    thumbnail = img.get('src')
            
            # Extract title dan link
            title_elem = li.find('h2')
            title = None
            href = None
            
            if title_elem:
                a_tag = title_elem.find('a')
                if a_tag:
                    title = a_tag.text.strip()
                    href = a_tag.get('href')
            else:
                # Fallback: cari a tag langsung
                a_tag = li.find('a', href=True)
                if a_tag and not a_tag.find('img'):  # Hindari link dari gambar
                    title = a_tag.text.strip()
                    href = a_tag.get('href')
            
            # Extract release date
            release_date = None
            date_div = li.find('div', class_='set')
            if date_div:
                release_date = date_div.text.strip()
            
            # Extract genre (jika ada)
            genres = []
            genre_div = li.find('div', class_='genrenya')
            if genre_div:
                genre_links = genre_div.find_all('a')
                genres = [g.text.strip() for g in genre_links]
            
            # Extract episode
            episode = None
            ep_div = li.find('div', class_='epz')
            if ep_div:
                episode = ep_div.text.strip()
            elif li.find('span', class_='ep'):
                episode = li.find('span', class_='ep').text.strip()
            
            # Extract rating
            rating = None
            rating_div = li.find('div', class_='rating')
            if rating_div:
                rating = rating_div.text.strip()
            
            if title and href:
                animes.append(Anime(
                    title=title,
                    original_page_url=href,
                    thumbnail=thumbnail,
                    release_date=release_date,
                    genre=genres,
                    rating=rating,
                    episodes=episode
                ))
        
        return animes
    
    def get_page_url(self, page=1, search_query=None, post_type='anime'):
        """Generate URL untuk halaman tertentu dengan dukungan search"""
        params = {}
        
        # Add post_type
        if post_type:
            params['post_type'] = post_type
        
        # Add search query
        if search_query:
            params['s'] = search_query
        
        if page <= 1:
            if params:
                query_string = urlencode(params)
                return f"{self.base_url}?{query_string}"
            return self.base_url
        else:
            # Halaman > 1: base_url/page/N/?params
            if params:
                query_string = urlencode(params)
                return f"{self.base_url}page/{page}/?{query_string}"
            return f"{self.base_url}page/{page}/"
    
    def scrape_page(self, page=1, search_query=None, post_type='anime'):
        """Scrape halaman tertentu dengan opsi search"""
        url = self.get_page_url(page, search_query, post_type)
        print(f"Scraping URL: {url}")
        
        try:
            response = self.session.get(url, timeout=self.timeout)
            response.raise_for_status()
            response.encoding = 'utf-8'
        except requests.exceptions.RequestException as e:
            print(f"Error fetching {url}: {e}")
            return None
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Extract animes
        animes = self._extract_animes(soup)
        
        # Extract pagination
        pagination = self._parse_pagination(soup)
        
        # Jika tidak ada pagination, buat default
        if not pagination:
            pagination = PaginationInfo(
                total_pages=1,
                current_page=page,
                has_next=False,
                has_previous=page > 1,
                next_page_url=self.get_page_url(page+1, search_query, post_type) if animes else None,
                previous_page_url=self.get_page_url(page-1, search_query, post_type) if page > 1 else None,
                page_numbers=[page],
                per_page=len(animes) or 20
            )
        
        return ScrapeResult(
            animes=animes,
            pagination=pagination,
            query=search_query
        )
    
    def scrape_range(self, start_page=1, end_page=None, search_query=None, post_type='anime'):
        """Scrape range halaman"""
        results = []
        current_page = start_page
        
        while True:
            print(f"Scraping page {current_page}...")
            result = self.scrape_page(current_page, search_query, post_type)
            
            if not result or not result.animes:
                print(f"No more animes found at page {current_page}")
                break
            
            results.append({
                "page": current_page,
                "animes": result.animes,
                "pagination": result.pagination,
                "query": result.query
            })
            
            print(f"  Found {len(result.animes)} animes")
            
            # Stop conditions
            if end_page and current_page >= end_page:
                print(f"Reached end page: {end_page}")
                break
            if not result.pagination or not result.pagination.has_next:
                print(f"No next page available")
                break
                
            current_page += 1
        
        return results
    
    def search(self, query, page=1):
        """Search animes dengan keyword"""
        return self.scrape_page(page=page, search_query=query, post_type='anime')
    
    def get_all_animes(self, max_pages=None):
        """Get semua animes (hati-hati, bisa ribuan)"""
        return self.scrape_range(start_page=1, end_page=max_pages, post_type='anime')
    
    def get_latest_animes(self, limit=20):
        """Get animes terbaru (halaman 1)"""
        result = self.scrape_page(page=1)
        if result:
            return result.animes[:limit]
        return []

class OtakuDesuScraperOnGoing:
    def __init__(self, base_url="https://otakudesu.best/"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.ua = UserAgent()
        self.timeout = 15
        self.retry_count = 3
    
    def _get_headers(self):
        return {
            'User-Agent': self.ua.random,
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'id,en-US;q=0.7,en;q=0.3',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Referer': self.base_url,
            'Upgrade-Insecure-Requests': '1',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'none',
            'Cache-Control': 'max-age=0',
        }
    
    def _make_request(self, url, max_retries=3):
        """Enhanced request dengan retry, delay, dan rotasi user-agent"""
        for attempt in range(max_retries):
            try:
                # Rotasi User-Agent each request
                headers = self._get_headers()
                
                # Random delay between 2-5 seconds
                time.sleep(random.uniform(2, 5))
                
                response = self.session.get(
                    url, 
                    headers=headers, 
                    timeout=self.timeout,
                    allow_redirects=True
                )
                
                if response.status_code == 200:
                    return response
                elif response.status_code == 403:
                    print(f"⚠️ Attempt {attempt + 1}: Blocked (403) for {url}")
                    
                    # Change IP with proxy (if have)
                    # self.session.proxies = {'http': 'http://new-proxy:port'}
                    
                    if attempt < max_retries - 1:
                        time.sleep(random.uniform(5, 10))
                        continue
                
                response.raise_for_status()
                
            except requests.exceptions.RequestException as e:
                print(f"⚠️ Attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(random.uniform(3, 7))
                else:
                    return None
        
        return None
    
    def get_ongoing_url(self, page=1):
        if page <= 1:
            return f"{self.base_url}/ongoing-anime/"
        return f"{self.base_url}/ongoing-anime/page/{page}/"
    
    def _extract_ongoing_anime(self, soup):
        ongoing_list = []
        
        venutama = soup.find('div', class_='venutama')
        if not venutama:
            print("❌ Cannot find div.venutama")
            return ongoing_list
        
        rseries = venutama.find('div', class_='rseries')
        if not rseries:
            print("❌ Cannot find div.rseries")
            return ongoing_list
        
        rapi = rseries.find('div', class_='rapi')
        if not rapi:
            print("❌ Cannot find div.rapi")
            return ongoing_list
        
        venz = rapi.find('div', class_='venz')
        if not venz:
            print("❌ Cannot find div.venz")
            return ongoing_list
        
        # Find UL
        ul = venz.find('ul')
        if not ul:
            print("❌ Cannot find ul")
            return ongoing_list
        
        # Find all LI directly under UL (not recursive)
        li_items = ul.find_all('li', recursive=False)
        print(f"✅ Found {len(li_items)} li items")
        
        for li in li_items:
            detpost = li.find('div', class_='detpost')
            if not detpost:
                continue
            
            # --- Episode (current episode) ---
            current_episode = None
            epz = detpost.find('div', class_='epz')
            if epz:
                episode_text = epz.text.strip()
                import re
                match = re.search(r'Episode\s*(\d+(?:\.\d+)?)', episode_text, re.IGNORECASE)
                if match:
                    current_episode = match.group(1)
                else:
                    current_episode = episode_text
            
            # --- Release Day (released_day) ---
            released_day = None
            epztipe = detpost.find('div', class_='epztipe')
            if epztipe:
                # Remove icon
                icon = epztipe.find('i')
                if icon:
                    icon.extract()
                released_day = epztipe.text.strip()
            
            # --- Release Date (release_date) ---
            release_date = None
            newnime = detpost.find('div', class_='newnime')
            if newnime:
                release_date = newnime.text.strip()
            
            # --- Thumbnail dan Title ---
            thumb_div = detpost.find('div', class_='thumb')
            if thumb_div:
                a_tag = thumb_div.find('a', href=True)
                if a_tag:
                    href = a_tag.get('href')
                    
                    # Thumbnail
                    thumbnail = None
                    thumbz = thumb_div.find('div', class_='thumbz')
                    if thumbz:
                        img = thumbz.find('img')
                        if img:
                            if img.get('srcset'):
                                thumbnail = img.get('srcset').split(',')[0].split(' ')[0]
                            elif img.get('src'):
                                thumbnail = img.get('src')
                    
                    # Title
                    title = None
                    h2_tag = thumbz.find('h2', class_='jdlflm') if thumbz else None
                    if h2_tag:
                        title = h2_tag.text.strip()
                    
                    if title and href:
                        # ✅ PERBAIKAN: Pakai field yang sesuai dengan dataclass Anime
                        ongoing_list.append(Anime(
                            title=title,
                            original_page_url=href,
                            thumbnail=thumbnail,
                            total_episodes=current_episode,  # ← Simpan episode terbaru di sini
                            released_day=released_day,        # ← Perbaiki nama field!
                            release_date=release_date,
                            status="Ongoing"                  # ← Set status otomatis
                        ))
        
        print(f"✅ Successfully extracted {len(ongoing_list)} ongoing anime")
        return ongoing_list
    
    def _parse_ongoing_pagination(self, soup):
        venutama = soup.find('div', class_='venutama')
        if not venutama:
            return None
        
        # Find pagination div
        pagination_div = venutama.find('div', class_='pagination')
        if not pagination_div:
            return None
        
        pagenavix = pagination_div.find('div', class_='pagenavix')
        if not pagenavix:
            return None
        
        # Current page
        current_page = 1
        current_span = pagenavix.find('span', class_='current', attrs={'aria-current': 'page'})
        if current_span and current_span.text.strip().isdigit():
            current_page = int(current_span.text.strip())
        
        # Total pages - find largest number in page links
        total_pages = 1
        page_links = pagenavix.find_all('a', class_='page-numbers')
        for link in page_links:
            if link.text.strip().isdigit():
                page_num = int(link.text.strip())
                if page_num > total_pages:
                    total_pages = page_num
        
        # Next URL
        next_page_url = None
        next_link = pagenavix.find('a', class_='next')
        if next_link and next_link.get('href'):
            next_page_url = next_link.get('href')
        
        # Prev URL
        previous_page_url = None
        prev_link = pagenavix.find('a', class_='prev')
        if prev_link and prev_link.get('href'):
            previous_page_url = prev_link.get('href')
        
        return PaginationInfo(
            total_pages=total_pages,
            current_page=current_page,
            has_next=next_page_url is not None,
            has_previous=previous_page_url is not None,
            next_page_url=next_page_url,
            previous_page_url=previous_page_url,
            page_numbers=[int(l.text) for l in page_links if l.text.strip().isdigit()]
        )
    
    def scrape_ongoing_page(self, page=1):
        url = self.get_ongoing_url(page)
        print(f"Scraping ongoing page: {url}")
        
        try:
            response = self.session.get(url, timeout=self.timeout)
            response.raise_for_status()
            response.encoding = 'utf-8'
        except requests.exceptions.RequestException as e:
            print(f"Error fetching {url}: {e}")
            return None
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Extract ongoing anime
        ongoing_list = self._extract_ongoing_anime(soup)
        
        # Extract pagination
        pagination = self._parse_ongoing_pagination(soup)
        
        # Fallback pagination
        if not pagination:
            pagination = PaginationInfo(
                total_pages=1,
                current_page=page,
                has_next=False,
                has_previous=page > 1,
                next_page_url=self.get_ongoing_url(page+1) if page < 10 else None,
                previous_page_url=self.get_ongoing_url(page-1) if page > 1 else None,
                page_numbers=[page]
            )
        
        return ScrapeResult(
            animes=ongoing_list,
            pagination=pagination
        )
        
    def scrape_ongoing_range(self, start_page=1, end_page=None):
        results = []
        current_page = start_page
        
        while True:
            print(f"Scraping ongoing page {current_page}...")
            result = self.scrape_ongoing_page(current_page)
            
            if not result or not result.animes:
                print(f"No more ongoing anime at page {current_page}")
                break
            
            results.append({
                "page": current_page,
                "animes": result.animes,  # ← Konsisten pakai 'animes'
                "pagination": result.pagination
            })
            
            print(f"  Found {len(result.animes)} ongoing anime")
            
            if end_page and current_page >= end_page:
                break
            if not result.pagination or not result.pagination.has_next:
                break
            
            current_page += 1
        
        return results

class OtakuDesuScraperGenreList:
    def __init__(self, base_url="https://otakudesu.best/"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.ua = UserAgent()
        self.timeout = 15
        self.retry_count = 3
    
    def _get_headers(self):
        return {
            'User-Agent': self.ua.random,
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'id,en-US;q=0.7,en;q=0.3',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Referer': self.base_url,
            'Upgrade-Insecure-Requests': '1',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'none',
            'Cache-Control': 'max-age=0',
        }
    
    def _make_request(self, url, max_retries=3):
        """Enhanced request dengan retry, delay, dan rotasi user-agent"""
        for attempt in range(max_retries):
            try:
                # Rotasi User-Agent each request
                headers = self._get_headers()
                
                # Random delay between 2-5 seconds
                time.sleep(random.uniform(2, 5))
                
                response = self.session.get(
                    url, 
                    headers=headers, 
                    timeout=self.timeout,
                    allow_redirects=True
                )
                
                if response.status_code == 200:
                    return response
                elif response.status_code == 403:
                    print(f"⚠️ Attempt {attempt + 1}: Blocked (403) for {url}")
                    
                    # Change IP with proxy (if have)
                    # self.session.proxies = {'http': 'http://new-proxy:port'}
                    
                    if attempt < max_retries - 1:
                        time.sleep(random.uniform(5, 10))
                        continue
                
                response.raise_for_status()
                
            except requests.exceptions.RequestException as e:
                print(f"⚠️ Attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(random.uniform(3, 7))
                else:
                    return None
        
        return None

    def get_genre_url(self):
       return f"{self.base_url}/genre-list/"
   
    def extract_genres(self, html_content):
        """Extract genre links from HTML content"""
        soup = BeautifulSoup(html_content, 'html.parser')
        genre_links = soup.find('ul', class_='genres').find_all('a')
        
        genres = []
        for link in genre_links:
            genre_name = link.text.strip()
            genre_url = link['href']
            
            genres.append(Genre(
                name=genre_name,
                ## merge base url + pathname genre
                url=self.base_url + genre_url.replace(self.base_url, "")
            ))
        
        return genres
   
    def scrape_genre_list(self):
        """Scrape genre list from the website"""
        url = self.get_genre_url()
        response = self._make_request(url)
        
        if not response:
            print(f"Failed to fetch genre list from {url}")
            return None
        
        return self.extract_genres(response.text)

class OtakuDesuEpisodeScraper:
    def __init__(self, base_url="https://otakudesu.best/"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.ua = UserAgent()
        self.timeout = 15
        self.retry_count = 3
    
    def _get_headers(self):
        return {
            'User-Agent': self.ua.random,
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'id,en-US;q=0.7,en;q=0.3',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Referer': self.base_url,
            'Upgrade-Insecure-Requests': '1',
        }
    
    def _make_request(self, url, max_retries=3):
        for attempt in range(max_retries):
            try:
                headers = self._get_headers()
                time.sleep(random.uniform(1, 3))
                
                response = self.session.get(
                    url, 
                    headers=headers, 
                    timeout=self.timeout,
                    allow_redirects=True
                )
                
                if response.status_code == 200:
                    return response
                elif response.status_code == 403:
                    print(f"⚠️ Attempt {attempt + 1}: Blocked (403) for {url}")
                    if attempt < max_retries - 1:
                        time.sleep(random.uniform(5, 10))
                        continue
                
                response.raise_for_status()
                
            except requests.exceptions.RequestException as e:
                print(f"⚠️ Attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(random.uniform(3, 7))
                else:
                    return None
        return None
    
    def get_anime_url(self, slug):
        """Generate URL detail anime dari slug"""
        return f"{self.base_url}/anime/{slug}"
    
    def extract_anime_detail(self, url):
        """Extract semua detail anime termasuk episode list"""
        print(f"🔍 Scraping anime detail: {url}")
        
        response = self._make_request(url)
        if not response:
            print("❌ Failed to fetch anime detail page")
            return None
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # ===== 1. TITLE =====
        title = None
        jdlrx = soup.find('div', class_='jdlrx')
        if jdlrx:
            h1 = jdlrx.find('h1')
            if h1:
                title = h1.text.strip()
        
        if not title:
            title = "Unknown Title"
        
        # ===== 2. THUMBNAIL =====
        thumbnail = None
        fotoanime = soup.find('div', class_='fotoanime')
        if fotoanime:
            img = fotoanime.find('img')
            if img:
                if img.get('srcset'):
                    thumbnail = img.get('srcset').split(',')[0].split(' ')[0]
                elif img.get('src'):
                    thumbnail = img.get('src')
        
        # ===== 3. INFO SECTION =====
        info_data = {
            'title_japanese': None,
            'score': None,
            'producer': None,
            'type': None,
            'status': None,
            'total_episodes': None,
            'duration': None,
            'release_date': None,
            'studio': None,
            'genre': [],
            'rating': None
        }
        
        info_section = soup.find('div', class_='infozingle')
        if info_section:
            # Cari semua <p> di dalam infozingle
            paragraphs = info_section.find_all('p')
            
            for p in paragraphs:
                span = p.find('span')
                if not span:
                    continue
                
                text = span.text.strip()
                
                # Parsing berdasarkan label
                if 'Judul' in text and ':' in text:
                    # Skip, ini sama dengan title utama
                    pass
                elif 'Japanese' in text and ':' in text:
                    info_data['title_japanese'] = text.replace('Japanese:', '').strip()
                elif 'Skor' in text and ':' in text:
                    info_data['score'] = text.replace('Skor:', '').strip()
                elif 'Produser' in text and ':' in text:
                    info_data['producer'] = text.replace('Produser:', '').strip()
                elif 'Tipe' in text and ':' in text:
                    info_data['type'] = text.replace('Tipe:', '').strip()
                elif 'Status' in text and ':' in text:
                    info_data['status'] = text.replace('Status:', '').strip()
                elif 'Total Episode' in text and ':' in text:
                    info_data['total_episodes'] = text.replace('Total Episode:', '').strip()
                elif 'Durasi' in text and ':' in text:
                    info_data['duration'] = text.replace('Durasi:', '').strip()
                elif 'Tanggal Rilis' in text and ':' in text:
                    info_data['release_date'] = text.replace('Tanggal Rilis:', '').strip()
                elif 'Studio' in text and ':' in text:
                    info_data['studio'] = text.replace('Studio:', '').strip()
                elif 'Genre' in text and ':' in text:
                    # Extract genre links
                    genre_links = span.find_all('a')
                    for link in genre_links:
                        info_data['genre'].append(Genre(
                            name=link.text.strip(),
                            url=link.get('href')
                        ))
        
        # ===== 4. EPISODE LIST =====
        episodes = []
        
        # Cari semua div.episodelist
        episode_lists = soup.find_all('div', class_='episodelist')
        
        for episode_list in episode_lists:
            # Cari ul di dalam episode list
            ul = episode_list.find('ul')
            if not ul:
                continue
            
            # Cari semua li
            li_items = ul.find_all('li')
            
            for li in li_items:
                # Cari link episode
                span = li.find('span')
                if not span:
                    continue
                
                a_tag = span.find('a')
                if not a_tag:
                    continue
                
                episode_title = a_tag.text.strip()
                episode_url = a_tag.get('href')
                
                # Cari tanggal rilis
                release_date = None
                zeebr = li.find('span', class_='zeebr')
                if zeebr:
                    release_date = zeebr.text.strip()
                
                episodes.append(Episode(
                    title=episode_title,
                    page_url=episode_url,
                    release_date=release_date
                ))
        
        # Urutkan episode dari yang terbaru ke terlama
        episodes.reverse()
        
        # ===== 5. BUILD ANIME OBJECT =====
        anime = Anime(
            title=title,
            title_japanese=info_data['title_japanese'],
            original_page_url=url,
            thumbnail=thumbnail,
            score=info_data['score'],
            producer=info_data['producer'],
            type=info_data['type'],
            status=info_data['status'],
            total_episodes=info_data['total_episodes'],
            duration=info_data['duration'],
            release_date=info_data['release_date'],
            studio=info_data['studio'],
            genre=info_data['genre'],
            rating=info_data['score'],  # Rating pakai score
            episodes=episodes
        )
        
        print(f"✅ Success: {title}")
        print(f"   Episodes: {len(episodes)}")
        print(f"   Status: {info_data['status']}")
        
        return anime
    
    def get_episode_by_number(self, slug, episode_number):
        """Get specific episode by number"""
        anime = self.get_anime_detail(slug)
        if not anime:
            return None
        
        for ep in anime.episodes:
            if f"Episode {episode_number}" in ep.title:
                return ep
        
        return None
    
    def get_latest_episodes(self, limit=5):
        """Get latest episodes from homepage (future implementation)"""
        # Ini bisa diimplementasi untuk scrape halaman utama
        pass

class OtakuDesuDetailEpisodeScraper:
    def __init__(self, base_url="https://otakudesu.best/"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.ua = UserAgent()
        self.timeout = 15
        self.retry_count = 3
    
    def _get_headers(self):
        return {
            'User-Agent': self.ua.random,
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'id,en-US;q=0.7,en;q=0.3',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Referer': self.base_url,
            'Upgrade-Insecure-Requests': '1',
        }
    
    def _make_request(self, url, max_retries=3):
        for attempt in range(max_retries):
            try:
                headers = self._get_headers()
                time.sleep(random.uniform(1, 3))
                
                response = self.session.get(
                    url, 
                    headers=headers, 
                    timeout=self.timeout,
                    allow_redirects=True
                )
                
                if response.status_code == 200:
                    return response
                elif response.status_code == 403:
                    print(f"⚠️ Attempt {attempt + 1}: Blocked (403) for {url}")
                    if attempt < max_retries - 1:
                        time.sleep(random.uniform(5, 10))
                        continue
                
                response.raise_for_status()
                
            except requests.exceptions.RequestException as e:
                print(f"⚠️ Attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(random.uniform(3, 7))
                else:
                    return None
        return None
    
    def _decode_mirror_data(self, data_content):
        """Decode base64 data untuk mirror links"""
        try:
            decoded = base64.b64decode(data_content).decode('utf-8')
            return json.loads(decoded)
        except:
            return None
    
    def extract_episode_info(self, url):
        print(f"🔍 Scraping episode page: {url}")
        response = self._make_request(url)
        if not response:
            print("❌ Failed to fetch episode page")
            return None
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # ===== 1. TITLE =====
        title_tag = soup.find('h1', class_='posttl')
        title = title_tag.text.strip() if title_tag else "No Title"
        
        # Extract episode number
        episode_number = None
        if title:
            match = re.search(r'Episode\s*(\d+(?:\.\d+)?)', title, re.IGNORECASE)
            if match:
                episode_number = match.group(1)
        
        # ===== 2. POSTED BY & RELEASE INFO =====
        posted_by = None
        release_time = None
        release_date = None
        
        kategoz = soup.find('div', class_='kategoz')
        if kategoz:
            # Posted by
            user_icon = kategoz.find('i', class_='fa fa-user')
            if user_icon:
                user_span = user_icon.find_next_sibling('span')
                if user_span:
                    posted_by = user_span.text.strip()
            
            # Release time
            clock_icon = kategoz.find('i', class_='fa fa-clock-o')
            if clock_icon:
                time_span = clock_icon.find_next_sibling('span')
                if time_span:
                    release_time = time_span.text.strip()
        
        # ===== 3. PLAYER URL =====
        player_url = None
        embed_holder = soup.find('div', id='embed_holder')
        if embed_holder:
            player_div = embed_holder.find('div', class_='responsive-embed-stream')
            if player_div:
                iframe = player_div.find('iframe', src=True)
                if iframe:
                    player_url = iframe['src']
        
        # ===== 4. PREVIOUS & NEXT EPISODE =====
        previous_episode_url = None
        next_episode_url = None
        see_all_episodes_url = None
        
        prevnext = soup.find('div', class_='prevnext')
        if prevnext:
            # Cari div.flir
            flir = prevnext.find('div', class_='flir')
            if flir:
                # Cari semua link di flir
                links = flir.find_all('a', href=True)
                
                for link in links:
                    href = link['href']
                    link_text = link.text.strip().lower()
                    link_title = link.get('title', '').lower()
                    
                    # Previous episode
                    if 'previous' in link_text or 'sebelumnya' in link_text or 'prev' in link_text:
                        previous_episode_url = href
                    # See all episodes
                    elif 'see all' in link_text or 'all episodes' in link_text or 'semua episode' in link_text:
                        see_all_episodes_url = href
            
            # Cek juga di dropdown untuk next/previous? 
            # Tidak ada next episode di halaman episode, hanya previous dan see all
        
        # ===== 5. LIST EPISODE LAINNYA (DROPDOWN) =====
        list_episode = []
        select_cog = soup.find('select', id='selectcog')
        if select_cog:
            options = select_cog.find_all('option')
            # Skip option pertama (Pilih Episode Lainnya)
            for option in options[1:]:
                ep_title = option.text.strip()
                ep_url = option.get('value')
                if ep_title and ep_url:
                    list_episode.append(ListOfEpisode(
                        name=ep_title,
                        page_url=ep_url
                    ))
        
        # ===== 6. DOWNLOAD LINKS =====
        download_links = []
        
        # Cari div.download
        download_div = soup.find('div', class_='download')
        if download_div:
            # Cari semua ul di dalam download
            uls = download_div.find_all('ul')
            
            for ul in uls:
                # Cari strong untuk quality & format
                strong_tag = ul.find('strong')
                if not strong_tag:
                    continue
                
                quality_text = strong_tag.text.strip()
                
                # Parse quality dan format (contoh: "Mp4 360p", "MKV 480p")
                format_quality = quality_text.split()
                file_format = format_quality[0] if len(format_quality) > 0 else "Unknown"
                quality = format_quality[1] if len(format_quality) > 1 else "Unknown"
                
                # Cari semua li
                lis = ul.find_all('li')
                for li in lis:
                    # Cari tag <i> untuk ukuran file
                    size_tag = li.find('i')
                    file_size = size_tag.text.strip() if size_tag else None
                    
                    # Cari semua link download
                    links = li.find_all('a', href=True)
                    for link in links:
                        href = link['href']
                        name = link.text.strip()
                        
                        if href and name and not href.startswith('#'):
                            download_links.append(DownloadLink(
                                name=f"{name} [{quality}]",
                                url=href,
                                size=file_size,
                                quality=quality,
                                format=file_format
                            ))
        
        episode = Episode(
            title=title,
            player_url=player_url,
            page_url=url,
            posted_by=posted_by,
            previous_episode_url=previous_episode_url,
            next_episode_url=next_episode_url,  # Tetap None karena memang tidak ada
            release_time=release_time,
            list_episode=list_episode,
            download_links=download_links,
            episode_number=episode_number
        )
        
        # Simpan see_all_episodes_url sebagai attribute tambahan jika perlu
        episode.see_all_episodes_url = see_all_episodes_url
        
        print(f"✅ Episode: {title}")
        print(f"   Previous: {previous_episode_url}")
        print(f"   See All: {see_all_episodes_url}")
        print(f"   Player: {player_url[:50] if player_url else 'None'}...")
        print(f"   Downloads: {len(download_links)} links")
        print(f"   Other episodes: {len(list_episode)}")
        
        return episode
    
def main():
    scraper = OtakuDesuScraper()
    scraperOngoing = OtakuDesuScraperOnGoing()
    scraperGenreList = OtakuDesuScraperGenreList()
    scraperDetail = OtakuDesuEpisodeScraper()
    scraperEpisodeInfo = OtakuDesuDetailEpisodeScraper()
    
    print("=== LATEST ANIMES (Page 1) ===")
    result_page1 = scraper.scrape_page(page=1)
    
    if result_page1:
        print(f"Animes found: {len(result_page1.animes)}")
        print(f"Pagination: Page {result_page1.pagination.current_page} of {result_page1.pagination.total_pages}")
        print(f"Has next: {result_page1.pagination.has_next}")
        
        # Convert to JSON
        output = {
            "success": True,
            "message": "Animes fetched successfully",
            "data": [dataclasses.asdict(m) for m in result_page1.animes],
            "pagination": dataclasses.asdict(result_page1.pagination)
        }
        
        # Save to file
        with open('results/otakudesu_latest.json', 'w', encoding='utf-8') as f:
            json.dump(output, f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_latest.json")
    
    print("\n" + "="*50 + "\n")
    
    print("=== SEARCH ANIMES ===")
    search_result = scraper.search(query="darwin", page=1)
    
    if search_result and search_result.animes:
        print(f"Search '{search_result.query}': Found {len(search_result.animes)} animes")
        
        for i, anime in enumerate(search_result.animes[:5], 1):
            print(f"  {i}. {anime.title}")
            print(f"     URL: {anime.original_page_url}")
            print(f"     Date: {anime.release_date}")
            print()
        
        # Save search result
        output = {
            "success": True,
            "message": "Animes fetched successfully",
            "query": search_result.query,
            "data": [dataclasses.asdict(m) for m in search_result.animes],
            "pagination": dataclasses.asdict(search_result.pagination)
        }
        
        with open('results/otakudesu_search.json', 'w', encoding='utf-8') as f:
            json.dump(output, f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_search.json")
    
    print("\n" + "="*50 + "\n")
    
    print("=== SCRAPE MULTIPLE PAGES (Limited to 3) ===")
    multiple_pages = scraper.scrape_range(start_page=1, end_page=3)
    
    if multiple_pages:
        total_animes = sum(len(page['animes']) for page in multiple_pages)
        print(f"Total {total_animes} animes from {len(multiple_pages)} pages")
        
        # Save all
        with open('results/otakudesu_multipage.json', 'w', encoding='utf-8') as f:
            json.dump(multiple_pages, f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_multipage.json")
    
    print("\n" + "="*50 + "\n")
    
    print("=== API RESPONSE FORMAT ===")
    if result_page1:
        api_response = {
            "status": 200,
            "success": True,
            "message": "Animes fetched successfully",
            "data": {
                "page": result_page1.pagination.current_page,
                "total_pages": result_page1.pagination.total_pages,
                "total_items": len(result_page1.animes),
                "has_next": result_page1.pagination.has_next,
                "has_previous": result_page1.pagination.has_previous,
                "next_page_url": result_page1.pagination.next_page_url,
                "previous_page_url": result_page1.pagination.previous_page_url,
                "animes": [dataclasses.asdict(m) for m in result_page1.animes]
            }
        }
        
        print(json.dumps(api_response, indent=2, ensure_ascii=False)[:500] + "...\n")


    print("=== ONGOING ANIME - PAGE 1 ===")
    result = scraperOngoing.scrape_ongoing_page(page=2)
    if result:
        print(f"Ongoing animes found: {len(result.animes)}")
        print(f"Pagination: Page {result.pagination.current_page} of {result.pagination.total_pages}")
        print(f"Has next: {result.pagination.has_next}")
        
        with open('results/otakudesu_ongoing.json', 'w', encoding='utf-8') as f:
            json.dump(dataclasses.asdict(result), f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_ongoing.json")

    print("\n" + "="*50 + "\n")
    
    print("=== GENRE LIST ===")
    genre_list = scraperGenreList.scrape_genre_list()
    if genre_list:
        print(f"Total genres found: {len(genre_list)}")
        for genre in genre_list:
            print(f"{genre.name} - {genre.url}")
        
        with open('results/otakudesu_genres.json', 'w', encoding='utf-8') as f:
            json.dump([dataclasses.asdict(m) for m in genre_list], f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_genres.json")


    targetUrl = "https://otakudesu.best/anime/oshi-ko-s3-sub-indo/"
    anime_detail = scraperDetail.extract_anime_detail(targetUrl)
    if anime_detail:
        print(f"Detail for {anime_detail.title}")
        print(f"URL: {anime_detail.original_page_url}")
        print(f"Release Date: {anime_detail.release_date}")
            
        with open('results/otakudesu_detail.json', 'w', encoding='utf-8') as f:
                json.dump(dataclasses.asdict(anime_detail), f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_detail.json")
    
    targetEpisodeUrl = "https://otakudesu.best/episode/onk-s3-episode-4-sub-indo/"
    anime_detail = scraperEpisodeInfo.extract_episode_info(targetEpisodeUrl)
    if anime_detail:
            
        with open('results/otakudesu_episode.json', 'w', encoding='utf-8') as f:
                json.dump(dataclasses.asdict(anime_detail), f, cls=EnhancedJSONEncoder, indent=2, ensure_ascii=False)
        print("Saved to results/otakudesu_episode.json")

if __name__ == "__main__":
    main()