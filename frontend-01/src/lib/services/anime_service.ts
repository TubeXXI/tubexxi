import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';



export class AnimeServiceImpl extends BaseService implements AnimeService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient) {
		super(event);
	}

	async GetLatest(page: number = 1): Promise<PaginatedResult<Anime>> {
		try {
			const response = await this.api.publicRequest<Anime[]>('GET', `/anime/latest?page=${page}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get latest animes');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination,
			}
		} catch (error) {
			console.error("Failed to get latest animes:", error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false,
				},
			};
		}
	}
	async Search(q: string, page: number = 1): Promise<PaginatedResult<Anime>> {
		try {
			const response = await this.api.publicRequest<Anime[]>('GET', `/anime/search?s=${q}&page=${page}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to search animes');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination,
			}
		} catch (error) {
			console.error("Failed to search animes:", error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false,
				},
			};
		}
	}
	async GetOngoing(page: number = 1): Promise<PaginatedResult<Anime>> {
		try {
			const response = await this.api.publicRequest<Anime[]>('GET', `/anime/ongoing?page=${page}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get ongoing animes');
			}
			return {
				data: response.data || [],
				pagination: response.pagination as ApiPagination,
			}
		} catch (error) {
			console.error("Failed to get ongoing animes:", error);
			return {
				data: [],
				pagination: {
					current_page: 0,
					limit: 0,
					total_items: 0,
					total_pages: 0,
					has_prev: false,
					has_next: false,
				},
			};
		}
	}
	async GetGenres(): Promise<AnimeGenre[]> {
		try {
			const response = await this.api.publicRequest<AnimeGenre[]>('GET', `/anime/genres`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get genres');
			}
			return response.data || [];
		} catch (error) {
			console.error("Failed to get genres:", error);
			return [];
		}
	}
	async GetDetail(url: string): Promise<Anime> {
		try {
			const response = await this.api.publicRequest<Anime>('GET', `/anime/detail?url=${url}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get anime detail');
			}
			return response.data || {} as Anime;
		} catch (error) {
			console.error("Failed to get anime detail:", error);
			return {} as Anime;
		}
	}
	async GetEpisode(url: string): Promise<AnimeEpisode[]> {
		try {
			const response = await this.api.publicRequest<AnimeEpisode[]>('GET', `/anime/episode?url=${url}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get anime episode');
			}
			return response.data || [];
		} catch (error) {
			console.error("Failed to get anime episode:", error);
			return [];
		}
	}
}
