import type { RequestEvent } from '@sveltejs/kit';
import { ApiClientHandler, AuthHelper, QueryHelper, LanguageHelper } from '@/helpers';
import {
	SettingServiceImpl,
	AuthServiceImpl,
	UserServiceImpl,
	AdminServiceImpl,
	MovieServiceImpl,
	SeriesServiceImpl,
	AnimeServiceImpl,
	ServerStatusServiceImpl,
	ApplicationServiceImpl,
} from '@/services';

export class Dependencies {
	public readonly apiClient: ApiClient;
	public readonly queryHelper: QueryHelper;
	public readonly languageHelper: LanguageHelper;
	public readonly authHelper: AuthHelper;

	public readonly settingService: SettingServiceImpl;
	public readonly authService: AuthServiceImpl;
	public readonly userService: UserServiceImpl;
	public readonly adminService: AdminServiceImpl;
	public readonly movieService: MovieServiceImpl;
	public readonly seriesService: SeriesServiceImpl;
	public readonly animeService: AnimeServiceImpl;
	public readonly serverStatusService: ServerStatusServiceImpl;
	public readonly applicationService: ApplicationServiceImpl;

	constructor(event: RequestEvent) {
		this.apiClient = new ApiClientHandler(event);
		this.queryHelper = new QueryHelper(event);
		this.languageHelper = new LanguageHelper(event);
		this.authHelper = new AuthHelper(event);

		this.settingService = new SettingServiceImpl(event, this.apiClient);
		this.authService = new AuthServiceImpl(event, this.apiClient);
		this.userService = new UserServiceImpl(event, this.apiClient);
		this.adminService = new AdminServiceImpl(event, this.apiClient);
		this.movieService = new MovieServiceImpl(event, this.apiClient);
		this.seriesService = new SeriesServiceImpl(event, this.apiClient);
		this.animeService = new AnimeServiceImpl(event, this.apiClient);
		this.serverStatusService = new ServerStatusServiceImpl(event, this.apiClient);
		this.applicationService = new ApplicationServiceImpl(event, this.apiClient);
	}
}
