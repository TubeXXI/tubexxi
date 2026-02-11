import re
from dataclasses import dataclass, field
from typing import List, Optional
from urllib.parse import urljoin

from bs4 import BeautifulSoup

from fetcher import fetch_html


@dataclass
class AnimeGenre:
    name: Optional[str] = None
    url: Optional[str] = None


@dataclass
class AnimeEpisodeItem:
    name: Optional[str] = None
    page_url: Optional[str] = None


@dataclass
class AnimeDownloadLink:
    name: Optional[str] = None
    url: Optional[str] = None
    size: Optional[str] = None
    quality: Optional[str] = None
    format: Optional[str] = None


@dataclass
class AnimeEpisode:
    title: Optional[str] = None
    player_url: Optional[str] = None
    page_url: Optional[str] = None
    posted_by: Optional[str] = None
    previous_episode_url: Optional[str] = None
    next_episode_url: Optional[str] = None
    see_all_episodes_url: Optional[str] = None
    release_date: Optional[str] = None
    release_time: Optional[str] = None
    episode_number: Optional[str] = None
    list_episode: List[AnimeEpisodeItem] = field(default_factory=list)
    download_links: List[AnimeDownloadLink] = field(default_factory=list)


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
    genre: List[AnimeGenre] = field(default_factory=list)
    rating: Optional[str] = None
    episodes: List[AnimeEpisode] = field(default_factory=list)


@dataclass
class PaginationAnime:
    total_pages: int = 1
    current_page: int = 1
    has_next: bool = False
    has_previous: bool = False
    next_page_url: Optional[str] = None
    previous_page_url: Optional[str] = None
    page_numbers: List[int] = field(default_factory=list)
    per_page: int = 20


@dataclass
class AnimeListResult:
    animes: List[Anime]
    pagination: PaginationAnime
    query: Optional[str] = None


class OtakuDesuScraper:
    def __init__(self, base_url: str = "https://otakudesu.best/") -> None:
        self.base_url = base_url.rstrip("/") + "/"

    def _abs_url(self, href: Optional[str]) -> Optional[str]:
        if not href:
            return None
        return urljoin(self.base_url, href)

    def _parse_pagination_archive(self, soup: BeautifulSoup, current_page_fallback: int = 1) -> PaginationAnime:
        pagination_div = soup.find("div", class_="pagination")
        if not pagination_div:
            return PaginationAnime(current_page=current_page_fallback, total_pages=current_page_fallback)

        total_pages = 1
        current_page = current_page_fallback

        naviright = pagination_div.find("span", class_="naviright")
        if naviright:
            page_text = naviright.get_text(" ", strip=True)
            m = re.search(r"Pages?\s*(\d+)\s*of\s*(\d+)", page_text, re.IGNORECASE)
            if m:
                current_page = int(m.group(1))
                total_pages = int(m.group(2))

        next_page_url = None
        previous_page_url = None
        page_numbers: List[int] = []

        navileft = pagination_div.find("span", class_="navileft")
        if navileft:
            next_link = navileft.find("a", string="»")
            if next_link and next_link.get("href"):
                next_page_url = self._abs_url(next_link.get("href"))
            prev_link = navileft.find("a", string="«")
            if prev_link and prev_link.get("href"):
                previous_page_url = self._abs_url(prev_link.get("href"))

            for a in navileft.find_all("a"):
                txt = a.get_text(strip=True)
                if txt.isdigit():
                    page_numbers.append(int(txt))

        if not next_page_url:
            link_next = soup.find("link", rel="next")
            if link_next and link_next.get("href"):
                next_page_url = self._abs_url(link_next.get("href"))

        if not previous_page_url:
            link_prev = soup.find("link", rel="prev")
            if link_prev and link_prev.get("href"):
                previous_page_url = self._abs_url(link_prev.get("href"))

        return PaginationAnime(
            total_pages=total_pages,
            current_page=current_page,
            has_next=next_page_url is not None,
            has_previous=previous_page_url is not None,
            next_page_url=next_page_url,
            previous_page_url=previous_page_url,
            page_numbers=sorted(set(page_numbers)),
        )

    def _extract_archive_animes(self, soup: BeautifulSoup) -> List[Anime]:
        animes: List[Anime] = []
        ul_item = soup.find("ul", class_="chivsrc")
        if not ul_item:
            return animes

        li_items = ul_item.find_all("li", recursive=False)
        for li in li_items:
            if li.find("div", class_="pagination"):
                continue

            thumbnail = None
            img = li.find("img")
            if img:
                if img.get("srcset"):
                    first = img.get("srcset").split(",")[0].strip().split(" ")[0]
                    thumbnail = first or None
                elif img.get("src"):
                    thumbnail = img.get("src")

            title = None
            href = None

            title_elem = li.find("h2")
            if title_elem:
                a_tag = title_elem.find("a", href=True)
                if a_tag:
                    title = a_tag.get_text(strip=True)
                    href = a_tag.get("href")
            else:
                a_tag = li.find("a", href=True)
                if a_tag and not a_tag.find("img"):
                    title = a_tag.get_text(strip=True)
                    href = a_tag.get("href")

            if not title or not href:
                continue

            release_date = None
            genres: List[AnimeGenre] = []

            date_div = li.find("div", class_="set")
            if date_div:
                date_text = date_div.get_text(" ", strip=True)
                if re.match(r"^Genres\s*:\s*", date_text, re.IGNORECASE):
                    names = re.sub(r"^Genres\s*:\s*", "", date_text, flags=re.IGNORECASE)
                    for name in [p.strip() for p in names.split(",") if p.strip()]:
                        genres.append(AnimeGenre(name=name, url=None))
                else:
                    release_date = date_text

            genre_div = li.find("div", class_="genrenya")
            if genre_div:
                for a in genre_div.find_all("a", href=True):
                    n = a.get_text(strip=True)
                    u = self._abs_url(a.get("href"))
                    if n:
                        genres.append(AnimeGenre(name=n, url=u))

            total_episodes = None
            ep_div = li.find("div", class_="epz")
            if ep_div:
                total_episodes = ep_div.get_text(strip=True)
            else:
                ep_span = li.find("span", class_="ep")
                if ep_span:
                    total_episodes = ep_span.get_text(strip=True)

            rating = None
            rating_div = li.find("div", class_="rating")
            if rating_div:
                rating = rating_div.get_text(strip=True)

            animes.append(Anime(
                title=title,
                original_page_url=self._abs_url(href),
                thumbnail=thumbnail,
                release_date=release_date,
                total_episodes=total_episodes,
                genre=genres,
                rating=rating,
            ))

        return animes

    def get_latest(self, page: int = 1) -> AnimeListResult:
        if page <= 1:
            url = f"{self.base_url}?post_type=anime"
        else:
            url = f"{self.base_url}page/{page}/?post_type=anime"

        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")
        animes = self._extract_archive_animes(soup)
        pagination = self._parse_pagination_archive(soup, current_page_fallback=page)
        pagination.per_page = len(animes) or pagination.per_page
        return AnimeListResult(animes=animes, pagination=pagination)

    def search(self, query: str, page: int = 1) -> AnimeListResult:
        query = query or ""
        if page <= 1:
            url = f"{self.base_url}?s={query}&post_type=anime"
        else:
            url = f"{self.base_url}?s={query}&post_type=anime&page={page}"

        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")
        animes = self._extract_archive_animes(soup)
        pagination = self._parse_pagination_archive(soup, current_page_fallback=page)
        pagination.per_page = len(animes) or pagination.per_page
        return AnimeListResult(animes=animes, pagination=pagination, query=query)

    def get_genres(self) -> List[AnimeGenre]:
        url = f"{self.base_url}genre-list/"
        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")
        genres: List[AnimeGenre] = []
        seen = set()

        for a in soup.find_all("a", href=True):
            href = a.get("href")
            if not href:
                continue
            if "/genres/" not in href:
                continue
            name = a.get_text(strip=True)
            abs_url = self._abs_url(href)
            if not abs_url or not name:
                continue
            key = (name.lower(), abs_url)
            if key in seen:
                continue
            seen.add(key)
            genres.append(AnimeGenre(name=name, url=abs_url))

        return genres

    def get_ongoing(self, page: int = 1) -> AnimeListResult:
        if page <= 1:
            url = f"{self.base_url}ongoing-anime/"
        else:
            url = f"{self.base_url}ongoing-anime/page/{page}/"

        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")

        animes: List[Anime] = []
        venutama = soup.find("div", class_="venutama")
        venz = None
        if venutama:
            rseries = venutama.find("div", class_="rseries")
            if rseries:
                rapi = rseries.find("div", class_="rapi")
                if rapi:
                    venz = rapi.find("div", class_="venz")

        if venz:
            ul = venz.find("ul")
            if ul:
                li_items = ul.find_all("li", recursive=False)
                for li in li_items:
                    detpost = li.find("div", class_="detpost")
                    if not detpost:
                        continue
                    current_episode = None
                    epz = detpost.find("div", class_="epz")
                    if epz:
                        episode_text = epz.get_text(strip=True)
                        m = re.search(r"Episode\s*(\d+(?:\.\d+)?)", episode_text, re.IGNORECASE)
                        current_episode = m.group(1) if m else episode_text

                    released_day = None
                    epztipe = detpost.find("div", class_="epztipe")
                    if epztipe:
                        icon = epztipe.find("i")
                        if icon:
                            icon.extract()
                        released_day = epztipe.get_text(strip=True)

                    release_date = None
                    newnime = detpost.find("div", class_="newnime")
                    if newnime:
                        release_date = newnime.get_text(strip=True)

                    thumb_div = detpost.find("div", class_="thumb")
                    if not thumb_div:
                        continue
                    a_tag = thumb_div.find("a", href=True)
                    if not a_tag:
                        continue
                    href = a_tag.get("href")

                    thumbnail = None
                    title = None
                    thumbz = thumb_div.find("div", class_="thumbz")
                    if thumbz:
                        img = thumbz.find("img")
                        if img:
                            if img.get("srcset"):
                                thumbnail = img.get("srcset").split(",")[0].split(" ")[0]
                            elif img.get("src"):
                                thumbnail = img.get("src")
                        h2_tag = thumbz.find("h2", class_="jdlflm")
                        if h2_tag:
                            title = h2_tag.get_text(strip=True)

                    if title and href:
                        animes.append(Anime(
                            title=title,
                            original_page_url=self._abs_url(href),
                            thumbnail=thumbnail,
                            status="Ongoing",
                            total_episodes=current_episode,
                            released_day=released_day,
                            release_date=release_date,
                        ))

        pagination = PaginationAnime(current_page=page, total_pages=page)
        if venutama:
            pagination_div = venutama.find("div", class_="pagination")
            if pagination_div:
                pagenavix = pagination_div.find("div", class_="pagenavix")
                if pagenavix:
                    current_span = pagenavix.find("span", class_=re.compile(r"\bcurrent\b"), attrs={"aria-current": "page"})
                    if current_span and current_span.get_text(strip=True).isdigit():
                        pagination.current_page = int(current_span.get_text(strip=True))
                    total_pages = 1
                    page_links = pagenavix.find_all("a", class_=re.compile(r"\bpage-numbers\b"))
                    for link in page_links:
                        if link.get_text(strip=True).isdigit():
                            total_pages = max(total_pages, int(link.get_text(strip=True)))
                    pagination.total_pages = total_pages

                    next_link = pagenavix.find("a", class_=re.compile(r"\bnext\b"), href=True)
                    if next_link:
                        pagination.next_page_url = self._abs_url(next_link.get("href"))
                        pagination.has_next = pagination.next_page_url is not None

                    prev_link = pagenavix.find("a", class_=re.compile(r"\bprev\b"), href=True)
                    if prev_link:
                        pagination.previous_page_url = self._abs_url(prev_link.get("href"))
                        pagination.has_previous = pagination.previous_page_url is not None

                    nums = []
                    for link in page_links:
                        t = link.get_text(strip=True)
                        if t.isdigit():
                            nums.append(int(t))
                    pagination.page_numbers = sorted(set(nums))

        pagination.per_page = len(animes) or pagination.per_page
        return AnimeListResult(animes=animes, pagination=pagination)

    def get_detail(self, url: str) -> Optional[Anime]:
        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")

        title = "Unknown Title"
        jdlrx = soup.find("div", class_="jdlrx")
        if jdlrx:
            h1 = jdlrx.find("h1")
            if h1:
                title = h1.get_text(strip=True) or title

        thumbnail = None
        fotoanime = soup.find("div", class_="fotoanime")
        if fotoanime:
            img = fotoanime.find("img")
            if img:
                if img.get("srcset"):
                    thumbnail = img.get("srcset").split(",")[0].split(" ")[0]
                elif img.get("src"):
                    thumbnail = img.get("src")

        info = {
            "title_japanese": None,
            "score": None,
            "producer": None,
            "type": None,
            "status": None,
            "total_episodes": None,
            "duration": None,
            "release_date": None,
            "studio": None,
            "genre": [],
        }

        info_section = soup.find("div", class_="infozingle")
        if info_section:
            for p in info_section.find_all("p"):
                span = p.find("span")
                if not span:
                    continue
                text = span.get_text(" ", strip=True)
                if "Japanese:" in text:
                    info["title_japanese"] = text.replace("Japanese:", "").strip() or None
                elif "Skor:" in text:
                    info["score"] = text.replace("Skor:", "").strip() or None
                elif "Produser:" in text:
                    info["producer"] = text.replace("Produser:", "").strip() or None
                elif "Tipe:" in text:
                    info["type"] = text.replace("Tipe:", "").strip() or None
                elif "Status:" in text:
                    info["status"] = text.replace("Status:", "").strip() or None
                elif "Total Episode:" in text:
                    info["total_episodes"] = text.replace("Total Episode:", "").strip() or None
                elif "Durasi:" in text:
                    info["duration"] = text.replace("Durasi:", "").strip() or None
                elif "Tanggal Rilis:" in text:
                    info["release_date"] = text.replace("Tanggal Rilis:", "").strip() or None
                elif "Studio:" in text:
                    info["studio"] = text.replace("Studio:", "").strip() or None
                elif "Genre" in text:
                    for link in span.find_all("a", href=True):
                        n = link.get_text(strip=True)
                        u = self._abs_url(link.get("href"))
                        if n:
                            info["genre"].append(AnimeGenre(name=n, url=u))

        episodes: List[AnimeEpisode] = []
        for episode_list in soup.find_all("div", class_="episodelist"):
            ul = episode_list.find("ul")
            if not ul:
                continue
            for li in ul.find_all("li"):
                span = li.find("span")
                if not span:
                    continue
                a_tag = span.find("a", href=True)
                if not a_tag:
                    continue
                episode_title = a_tag.get_text(strip=True)
                episode_url = a_tag.get("href")
                if not episode_title or not episode_url:
                    continue
                release_date = None
                zeebr = li.find("span", class_="zeebr")
                if zeebr:
                    release_date = zeebr.get_text(strip=True)
                episodes.append(AnimeEpisode(
                    title=episode_title,
                    page_url=self._abs_url(episode_url),
                    release_date=release_date,
                ))
        episodes.reverse()

        return Anime(
            title=title,
            title_japanese=info["title_japanese"],
            original_page_url=url,
            thumbnail=thumbnail,
            score=info["score"],
            producer=info["producer"],
            type=info["type"],
            status=info["status"],
            total_episodes=info["total_episodes"],
            duration=info["duration"],
            release_date=info["release_date"],
            studio=info["studio"],
            genre=info["genre"],
            rating=info["score"],
            episodes=episodes,
        )

    def get_episode(self, url: str) -> Optional[AnimeEpisode]:
        html = fetch_html(url)
        soup = BeautifulSoup(html, "html.parser")

        title_tag = soup.find("h1", class_="posttl")
        title = title_tag.get_text(strip=True) if title_tag else "No Title"

        episode_number = None
        m = re.search(r"Episode\s*(\d+(?:\.\d+)?)", title or "", re.IGNORECASE)
        if m:
            episode_number = m.group(1)

        posted_by = None
        release_time = None
        kategoz = soup.find("div", class_="kategoz")
        if kategoz:
            user_icon = kategoz.find("i", class_="fa fa-user")
            if user_icon:
                user_span = user_icon.find_next_sibling("span")
                if user_span:
                    posted_by = user_span.get_text(strip=True)
            clock_icon = kategoz.find("i", class_="fa fa-clock-o")
            if clock_icon:
                time_span = clock_icon.find_next_sibling("span")
                if time_span:
                    release_time = time_span.get_text(strip=True)

        player_url = None
        embed_holder = soup.find("div", id="embed_holder")
        if embed_holder:
            player_div = embed_holder.find("div", class_="responsive-embed-stream")
            if player_div:
                iframe = player_div.find("iframe", src=True)
                if iframe:
                    player_url = iframe.get("src")

        previous_episode_url = None
        next_episode_url = None
        see_all_episodes_url = None

        prevnext = soup.find("div", class_="prevnext")
        if prevnext:
            flir = prevnext.find("div", class_="flir")
            if flir:
                for link in flir.find_all("a", href=True):
                    href = link.get("href")
                    link_text = link.get_text(strip=True).lower()
                    if "previous" in link_text or "sebelumnya" in link_text or "prev" in link_text:
                        previous_episode_url = self._abs_url(href)
                    elif "see all" in link_text or "all episodes" in link_text or "semua episode" in link_text:
                        see_all_episodes_url = self._abs_url(href)

        list_episode: List[AnimeEpisodeItem] = []
        select_cog = soup.find("select", id="selectcog")
        if select_cog:
            options = select_cog.find_all("option")
            for option in options[1:]:
                ep_title = option.get_text(strip=True)
                ep_url = option.get("value")
                if ep_title and ep_url:
                    list_episode.append(AnimeEpisodeItem(name=ep_title, page_url=self._abs_url(ep_url)))

        download_links: List[AnimeDownloadLink] = []
        download_div = soup.find("div", class_="download")
        if download_div:
            for ul in download_div.find_all("ul"):
                strong_tag = ul.find("strong")
                if not strong_tag:
                    continue
                quality_text = strong_tag.get_text(strip=True)
                parts = quality_text.split()
                file_format = parts[0] if len(parts) > 0 else None
                quality = parts[1] if len(parts) > 1 else None

                for li in ul.find_all("li"):
                    size_tag = li.find("i")
                    file_size = size_tag.get_text(strip=True) if size_tag else None
                    for link in li.find_all("a", href=True):
                        href = link.get("href")
                        name = link.get_text(strip=True)
                        if not href or not name or href.startswith("#"):
                            continue
                        display_name = name
                        if quality:
                            display_name = f"{name} [{quality}]"
                        download_links.append(AnimeDownloadLink(
                            name=display_name,
                            url=self._abs_url(href),
                            size=file_size,
                            quality=quality,
                            format=file_format,
                        ))

        return AnimeEpisode(
            title=title,
            player_url=player_url,
            page_url=url,
            posted_by=posted_by,
            previous_episode_url=previous_episode_url,
            next_episode_url=next_episode_url,
            see_all_episodes_url=see_all_episodes_url,
            release_time=release_time,
            list_episode=list_episode,
            download_links=download_links,
            episode_number=episode_number,
        )

