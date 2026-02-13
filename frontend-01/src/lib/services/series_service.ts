import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';

export class SeriesServiceImpl extends BaseService implements SeriesService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}
	async GetSeriesHome(): Promise<HomeScrapperResponse | Error> {
		try {
			const response = await this.api.publicRequest<HomeScrapperResponse>('GET', '/series/home');
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
	async GetSeriesByGenre(slug: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/genre/${slug}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series by genre');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get series by genre:', error);
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
	async GetSeriesByCountry(country: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/country/${country}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series by country');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get series by country:', error);
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
	async GetSeriesByYear(year: number, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/year/${year}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series by year');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get series by year:', error);
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
	async SearchSeries(query: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/search?s=${query}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to search series');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to search series:', error);
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
	async GetSeriesByFeature(type: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/featured/${type}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series by feature');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get series by feature:', error);
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
	async GetSeriesSpecialPage(path: string, page: number = 1): Promise<PaginatedResult<Movie>> {
		try {
			const response = await this.api.publicRequest<Movie[]>(
				'GET',
				`/series/special/${path}&page=${page}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get special page');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination
			};
		} catch (error) {
			console.error('Failed to get series special page:', error);
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
	async GetSeriesDetail(slug: string): Promise<SeriesDetail | null> {
		try {
			const response = await this.api.publicRequest<SeriesDetail>('GET', `/series/detail/${slug}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series detail');
			}
			return response.data || null;
		} catch (error) {
			console.error('Failed to get series detail:', error);
			return null;
		}
	}
	async GetSeriesEpisode(url: string): Promise<SeriesEpisode | null> {
		try {
			const response = await this.api.publicRequest<SeriesEpisode>(
				'GET',
				`/series/episode?url=${url}`
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get series episode');
			}
			return response.data || null;
		} catch (error) {
			console.error('Failed to get series episode:', error);
			return null;
		}
	}
}
