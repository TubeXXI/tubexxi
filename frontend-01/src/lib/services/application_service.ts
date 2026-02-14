import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';

export class ApplicationServiceImpl extends BaseService implements ApplicationService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}
	async RegisterApplication(
		apps: { package_name: string; key: string; value: string; description?: string; group_name: string }[]
	): Promise<boolean> {
		try {
			const endpoint = `/applications/protected`;

			const response = await this.api.authRequest('POST', endpoint, apps);
			if (!response.success) {
				throw new Error(response.message || 'Failed to register application');
			}
			return response.success;
		} catch (error) {
			console.error('Error registering application:', error);
			return false;
		}
	}
	async UpdateApplication(
		package_name: string,
		apps: { key: string; value: string; description?: string; group_name: string }[]
	): Promise<boolean> {
		try {
			const endpoint = `/applications/protected/${package_name}`;

			const response = await this.api.authRequest('PUT', endpoint, apps);
			if (!response.success) {
				throw new Error(response.message || 'Failed to update application');
			}
			return response.success;
		} catch (error) {
			console.error('Error updating application:', error);
			return false;
		}
	}
	async Search(query: QueryParams): Promise<PaginatedResult<ApplicationResponse>> {
		try {
			const params = new URLSearchParams();
			Object.entries(query).forEach(([key, value]) => {
				if (value !== undefined && value !== null && value !== '') {
					params.append(key, String(value));
				}
			});
			const queryString = params.toString();
			const endpoint = `/applications/protected/search${queryString ? `?${queryString}` : ''}`;

			const response = await this.api.authRequest('GET', endpoint);
			// console.log(response);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get applications');
			}
			return {
				data: response.data as ApplicationResponse[] || [],
				pagination: response.pagination as ApiPagination,
			}
		} catch (error) {
			console.error('Error fetching applications:', error);
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
	async GetByPackageName(packageName: string): Promise<ApplicationResponse | null> {
		try {
			const endpoint = `/applications/protected/${packageName}`;

			const response = await this.api.authRequest('GET', endpoint);
			if (!response.success) {
				throw new Error(response.message || 'Failed to get application');
			}
			return response.data as ApplicationResponse
		} catch (error) {
			console.error('Error fetching application:', error);
			return null;
		}
	}
	async Delete(packageName: string): Promise<boolean> {
		try {
			const endpoint = `/applications/protected/${packageName}`;

			const response = await this.api.authRequest('DELETE', endpoint);
			if (!response.success) {
				throw new Error(response.message || 'Failed to delete application');
			}
			return response.success;
		} catch (error) {
			console.error('Error deleting application:', error);
			return false;
		}
	}
	async BulkDelete(packageNames: string[]): Promise<boolean> {
		try {
			const endpoint = `/applications/protected/bulk`;

			const response = await this.api.authRequest('DELETE', endpoint, {
				package_names: packageNames,
			});
			if (!response.success) {
				throw new Error(response.message || 'Failed to delete applications');
			}
			return response.success;
		} catch (error) {
			console.error('Error deleting applications:', error);
			return false;
		}
	}
}
