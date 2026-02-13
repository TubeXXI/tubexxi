import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';

export class MovieServiceImpl extends BaseService implements MovieService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}

	async GetHome(): Promise<HomeScrapperResponse | Error> {
		try {
			const response = await this.api.publicRequest<HomeScrapperResponse>('GET', '/movies/home');
			if (!response.success) {
				throw new Error(response.message || 'Failed to get home scrapper');
			}
			return (
				response.data ||
				({
					key: null,
					value: null,
					view_all_url: null
				} as HomeScrapperResponse)
			);
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async GetMoviesByGenre(slug: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/genre/${slug}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get movies by genre');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get movies by genre:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async GetMoviesByCountry(country: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/country/${country}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get movies by country');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get movies by country:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async GetMoviesByYear(year: number, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/year/${year}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get movies by year');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get movies by year:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async SearchMovies(query: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/search?s=${query}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to search movies');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to search movies:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async GetMoviesByFeature(type: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/featured/${type}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get movies by feature');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get movies by feature:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async GetSpecialPage(path: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/movies/special/${path}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get special page');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get special page:', error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false
				}
			};
		}
	}
	async GetMovieDetail(slug: string): Promise<MovieDetail | null> {
		try {
			const response = await this.api.publicRequest<MovieDetail>('GET', `/movies/detail/${slug}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get movie detail');
			}
			return response.data || null;
		} catch (error) {
			console.error('Failed to get movie detail:', error);
			return null;
		}
	}
}
