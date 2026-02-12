import type { RequestEvent } from '@sveltejs/kit';
import { ApiClientHandler, AuthHelper, QueryHelper, LanguageHelper } from '@/helpers';

export class Dependencies {
	public readonly apiClient: ApiClient;
	public readonly queryHelper: QueryHelper;
	public readonly languageHelper: LanguageHelper;
	public readonly authHelper: AuthHelper;

	constructor(event: RequestEvent) {
		this.apiClient = new ApiClientHandler(event);
		this.queryHelper = new QueryHelper(event);
		this.languageHelper = new LanguageHelper(event);
		this.authHelper = new AuthHelper(event);
	}
}
