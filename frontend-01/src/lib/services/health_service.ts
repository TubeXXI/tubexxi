import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';


export class ServerStatusServiceImpl extends BaseService implements ServerStatusService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}
	async GetServerHealth(): Promise<ServerHealthResponse | null> {
		try {
			const response = await this.api.authRequest('GET', '/health/protected/check');
			if (!response.success) {
				throw new Error(response.message || 'Failed to get server health');
			}
			return response.data as ServerHealthResponse || {};
		} catch (error) {
			console.error('Error fetching server health:', error);
			return null;
		}
	}
	async GetServerLogs(page: number = 1, limit: number = 50): Promise<PaginatedResult<ServerLogsResponse> | null> {
		try {
			const response = await this.api.authRequest('GET', `/health/protected/log?page=${page}&limit=${limit}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get server logs');
			}
			return {
				data: response.data as ServerLogsResponse[] || [],
				pagination: response.pagination as ApiPagination,
			};
		} catch (error) {
			console.error('Error fetching server logs:', error);
			return null;
		}
	}
	async ClearServerLogs(): Promise<void | Error> {
		try {
			const response = await this.api.authRequest('POST', '/health/protected/log');
			if (!response.success) {
				throw new Error(response.message || 'Failed to clear server logs');
			}
		} catch (error) {
			console.error('Error clearing server logs:', error);
			return error instanceof Error ? error : new Error('Failed to clear server logs');
		}
	}
}
