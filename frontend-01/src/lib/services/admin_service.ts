import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';
import type { SetRoleSchema } from '@/utils/schema';

export class AdminServiceImpl extends BaseService implements AdminService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}
	async SetRole(data: SetRoleSchema): Promise<string | Error> {
		try {
			const response = await this.api.authRequest('POST', '/admin/users/protected/set-role', data);
			if (!response.success) {
				throw new Error(response.message || 'Failed to set role');
			}
			return response.message || 'Role set successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async SearchUser(query: QueryParams): Promise<PaginatedResult<User>> {
		try {
			const params = new URLSearchParams();
			Object.entries(query).forEach(([key, value]) => {
				if (value !== undefined && value !== null && value !== '') {
					params.append(key, String(value));
				}
			});
			const queryString = params.toString();
			const endpoint = `/admin/users/protected/search${queryString ? `?${queryString}` : ''}`;

			const response = await this.api.authRequest('GET', endpoint);
			// console.log(response);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get users');
			}
			return {
				data: response.data as User[] || [],
				pagination: response.pagination as ApiPagination,
			}
		} catch (error) {
			console.error('Error fetching users:', error);
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
	async HardDeleteUser(id: string): Promise<string | Error> {
		try {
			const response = await this.api.authRequest('DELETE', `/admin/users/protected/${id}`);
			if (!response.success) {
				throw new Error(response.message || 'Failed to delete user');
			}
			return response.message || 'User deleted successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async BulkDeleteUser(ids: string[]): Promise<string | Error> {
		try {
			const response = await this.api.authRequest('DELETE', '/admin/users/protected/bulk', { ids });
			if (!response.success) {
				throw new Error(response.message || 'Failed to delete users');
			}
			return response.message || 'Users deleted successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
}
