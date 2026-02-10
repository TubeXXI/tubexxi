from concurrent import futures
import logging
import grpc
import scraper_pb2
import scraper_pb2_grpc
from scraper import HomeScraper, MovieListScraper
from fetcher import fetch_html

class ScraperService(scraper_pb2_grpc.ScraperServiceServicer):
    def ScrapeHome(self, request, context):
        url = "https://tv8.lk21official.cc/"
        try:
            html_content = fetch_html(url)
            scraper = HomeScraper(html_content, base_url=url)
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

    def ScrapeList(self, request, context):
        url = request.url
        if not url:
             context.set_details("URL is required")
             context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
             return scraper_pb2.ListResponse()

        try:
            html_content = fetch_html(url)
            # Determine base_url from request url or default
            # Ideally base_url should be the domain
            base_url = "https://tv8.lk21official.cc"
            
            scraper = MovieListScraper(html_content, base_url=base_url)
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
