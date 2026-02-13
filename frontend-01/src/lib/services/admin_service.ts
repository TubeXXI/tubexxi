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
}
