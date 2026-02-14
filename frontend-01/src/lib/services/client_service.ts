import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';
import type { ContactSchema } from '@/utils/schema';


export class ClientServiceImpl extends BaseService implements ClientService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}

	async SendContact(data: ContactSchema): Promise<void | Error> {
		try {
			const response = await this.api.publicRequest('POST', '/client/web/contact', data);
			if (!response.success) {
				throw new Error(response.message || 'Failed to send contact');
			}
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown error');
		}
	}
}