import json
import os
import dataclasses
import argparse
from scraper import HomeScraper, MovieListScraper
from fetcher import fetch_html
try:
    from dotenv import load_dotenv
except ImportError:
    # Fallback: define a no-op load_dotenv if python-dotenv is not installed
    def load_dotenv(*args, **kwargs):
        pass

load_dotenv()

class EnhancedJSONEncoder(json.JSONEncoder):
        def default(self, o):
            if dataclasses.is_dataclass(o):
                return dataclasses.asdict(o)
            return super().default(o)

def main():
    parser = argparse.ArgumentParser(description='Movie Scraper')
    parser.add_argument('--type', type=str, default='home', choices=['home', 'list'], help='Scrape type: home or list')
    parser.add_argument('--url', type=str, help='URL to scrape (required for list type)')
    parser.add_argument('--local', action='store_true', help='Use local example.html')
    args = parser.parse_args()

    html_content = ""

    if args.local:
        file_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'docs', 'example.html')
        if not os.path.exists(file_path):
            print(f"Error: File not found at {file_path}")
            return
        with open(file_path, 'r', encoding='utf-8') as f:
            html_content = f.read()
    else:
        target_url = args.url if args.url else os.getenv("MOVIE_BASE_URL")
        try:
            html_content = fetch_html(target_url)
        except Exception as e:
            print(f"Failed to fetch content from {target_url}")
            return

    try:
        if args.type == 'home':
            scraper = HomeScraper(html_content)
            results = scraper.scrape()
        else:
            scraper = MovieListScraper(html_content)
            results = scraper.scrape()
        
        print(json.dumps(results, cls=EnhancedJSONEncoder, indent=2))
        
    except Exception as e:
        print(f"An error occurred during scraping: {e}")

if __name__ == "__main__":
    main()
