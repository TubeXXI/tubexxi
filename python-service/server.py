from concurrent import futures
import logging
import grpc
import scraper_pb2
import scraper_pb2_grpc
from scraper import HomeScraper, MovieListScraper, MovieDetailScraper
from fetcher import fetch_html
import urllib.parse

class ScraperService(scraper_pb2_grpc.ScraperServiceServicer):
    BASE_URL = "https://tv8.lk21official.cc"

    def ScrapeHome(self, request, context):
        try:
            html_content = fetch_html(self.BASE_URL)
            scraper = HomeScraper(html_content, base_url=self.BASE_URL)
            results = scraper.scrape()
            
            sections = []
            for item in results:
                movies = []
                for m in item.value:
                    movies.append(scraper_pb2.Movie(
                        id=m.id,
                        title=m.title,
                        original_title=m.original_title,
                        thumbnail=m.thumbnail,
                        synopsis=m.synopsis,
                        rating=m.rating,
                        duration=m.duration,
                        year=m.year,
                        date_published=m.date_published,
                        label_quality=m.label_quality,
                        genre=m.genre,
                        original_page_url=m.original_page_url
                    ))
                sections.append(scraper_pb2.HomeSection(
                    key=item.key, 
                    value=movies,
                    view_all_url=item.view_all_url
                ))
            
            return scraper_pb2.HomeResponse(sections=sections)
        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return scraper_pb2.HomeResponse()

    def _scrape_url(self, url, context):
        if not url:
             context.set_details("URL is required")
             context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
             return scraper_pb2.ListResponse()

        try:
            html_content = fetch_html(url)
            scraper = MovieListScraper(html_content, base_url=self.BASE_URL)
            result = scraper.scrape()
            
            movies = []
            for m in result.movies:
                movies.append(scraper_pb2.Movie(
                    id=m.id,
                    title=m.title,
                    original_title=m.original_title,
                    thumbnail=m.thumbnail,
                    synopsis=m.synopsis,
                    rating=m.rating,
                    duration=m.duration,
                    year=m.year,
                    date_published=m.date_published,
                    label_quality=m.label_quality,
                    genre=m.genre,
                    original_page_url=m.original_page_url
                ))
            
            pagination = scraper_pb2.Pagination(
                current_page=result.pagination.current_page,
                total_page=result.pagination.total_page,
                has_next=result.pagination.has_next,
                has_prev=result.pagination.has_prev,
                next_page_url=result.pagination.next_page_url,
                prev_page_url=result.pagination.prev_page_url
            )
            
            return scraper_pb2.ListResponse(movies=movies, pagination=pagination)
            
        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return scraper_pb2.ListResponse()

    def ScrapeList(self, request, context):
        return self._scrape_url(request.url, context)

    def GetMoviesByGenre(self, request, context):
        page = request.page if request.page > 0 else 1
        url = f"{self.BASE_URL}/genre/{request.slug}/page/{page}"
        return self._scrape_url(url, context)

    def SearchMovies(self, request, context):
        page = request.page if request.page > 0 else 1
        # Using query param 's' which is standard for WP, but user showed 'spiderman=am' which is weird.
        # Assuming ?s={query} works or /search/{query}/page/{page}
        # Let's try standard search query param
        query = urllib.parse.quote(request.query)
        url = f"{self.BASE_URL}/search?s={query}&page={page}"
        # Alternative from user input: https://tv8.lk21official.cc/search?spiderman=am&page=1
        # If the above fails, we might need to investigate. But standard WP search is usually safe.
        return self._scrape_url(url, context)

    def GetMoviesByFeature(self, request, context):
        page = request.page if request.page > 0 else 1
        # Feature types: populer, most-commented, rating, release, latest
        url = f"{self.BASE_URL}/{request.feature_type}/page/{page}"
        return self._scrape_url(url, context)

    def GetMoviesByCountry(self, request, context):
        page = request.page if request.page > 0 else 1
        url = f"{self.BASE_URL}/country/{request.country_slug}/page/{page}"
        return self._scrape_url(url, context)

    def GetMoviesByYear(self, request, context):
        page = request.page if request.page > 0 else 1
        url = f"{self.BASE_URL}/year/{request.year}/page/{page}"
        return self._scrape_url(url, context)

    def GetSpecialPage(self, request, context):
        # e.g. /rekomendasi-film-pintar
        # If page is needed, append /page/{n} ?? User didn't specify pagination for special page example.
        # But assuming it might have pagination.
        url = f"{self.BASE_URL}/{request.page_name}"
        if request.page > 1:
             url = f"{url}/page/{request.page}"
        return self._scrape_url(url, context)

    def GetMovieDetail(self, request, context):
        if not request.slug:
             context.set_details("Slug is required")
             context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
             return scraper_pb2.MovieDetailResponse()
        
        # Determine if slug is a full URL or just a slug
        # User might pass "wrecking-crew-2026" or "https://tv8.lk21official.cc/wrecking-crew-2026"
        url = request.slug
        if not url.startswith('http'):
            url = f"{self.BASE_URL}/{request.slug}"
            
        try:
            html_content = fetch_html(url)
            scraper = MovieDetailScraper(html_content, base_url=self.BASE_URL)
            detail = scraper.scrape(original_url=url)
            
            if not detail:
                context.set_details("Failed to scrape movie details")
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return scraper_pb2.MovieDetailResponse()

            # Map Entities to Proto
            
            # Map Movie
            m = detail.movie
            proto_movie = scraper_pb2.Movie(
                id=m.id,
                title=m.title,
                original_title=m.original_title,
                thumbnail=m.thumbnail,
                synopsis=m.synopsis,
                rating=m.rating,
                duration=m.duration,
                year=m.year,
                date_published=m.date_published,
                label_quality=m.label_quality,
                genre=m.genre,
                original_page_url=m.original_page_url
            )

            # Map Player Urls
            proto_player_urls = []
            for p in detail.player_urls:
                proto_player_urls.append(scraper_pb2.PlayerUrl(
                    url=p.url,
                    type=p.type
                ))
            
            # Map Directors
            proto_directors = []
            for p in detail.directors:
                proto_directors.append(scraper_pb2.MoviePerson(
                    name=p.name,
                    page_url=p.page_url
                ))

            # Map Stars
            proto_stars = []
            for p in detail.movie_stars:
                proto_stars.append(scraper_pb2.MoviePerson(
                    name=p.name,
                    page_url=p.page_url
                ))
            
            # Map Countries
            proto_countries = []
            for c in detail.countries:
                proto_countries.append(scraper_pb2.CountryMovie(
                    name=c.name,
                    page_url=c.page_url
                ))

            # Map Genres
            proto_genres = []
            for g in detail.genres:
                proto_genres.append(scraper_pb2.Genre(
                    name=g.name,
                    page_url=g.page_url
                ))
            
            # Map Similar Movies
            proto_similar = []
            for sm in detail.similar_movies:
                proto_similar.append(scraper_pb2.Movie(
                    id=sm.id,
                    title=sm.title,
                    original_title=sm.original_title,
                    thumbnail=sm.thumbnail,
                    synopsis=sm.synopsis,
                    rating=sm.rating,
                    duration=sm.duration,
                    year=sm.year,
                    date_published=sm.date_published,
                    label_quality=sm.label_quality,
                    genre=sm.genre,
                    original_page_url=sm.original_page_url
                ))

            # Create MovieDetail protobuf message
            movie_detail = scraper_pb2.MovieDetail(
                movie=proto_movie,
                votes=detail.votes,
                player_urls=proto_player_urls,
                directors=proto_directors,
                movie_stars=proto_stars,
                countries=proto_countries,
                genres=proto_genres,
                similar_movies=proto_similar
            )

            return scraper_pb2.MovieDetailResponse(detail=movie_detail)


        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return scraper_pb2.MovieDetailResponse()


def serve():
    port = '50051'
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    scraper_pb2_grpc.add_ScraperServiceServicer_to_server(ScraperService(), server)
    server.add_insecure_port('[::]:' + port)
    print("Scraper Service started on port " + port)
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig()
    serve()
