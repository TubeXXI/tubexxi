import type {
	RegisterSchema,
	ResetPasswordSchema,
	VerifyEmailSchema,
	ChangePasswordSchema,
	SetRoleSchema,
	UpdateProfileSchema,
	UpdatePasswordSchema,
	ContactSchema,
	RegisterAppSchema,
	UpdateApplicationSchema,
	UpdatePlatformSchema,
	DownloadVideoSchema,
	WebErrorReportSchema
} from '@/utils/schema';

declare global {
	interface Window {
		gc: NodeJS.GCFunction | undefined;
		adsbygoogle: any[] | undefined;
	}

	// ==========================================
	// Database Models (Matches 000001_init_schema.up.sql)
	// ==========================================

	type Role = {
		id: string;
		name: string;
		slug: string;
		description?: string | null;
		level: number;
		created_at: string;
		updated_at: string;
		deleted_at?: Date | null;
	};

	type User = {
		id: string;
		email: string;
		full_name: string;
		phone?: string | null;
		avatar_url: string;
		role_id: string;
		two_fa_secret?: string | null;
		is_active: boolean;
		is_verified: boolean;
		email_verified_at?: Date | string | null;
		last_login_at?: Date | string | null;
		created_at: string;
		updated_at: string;
		deleted_at?: string | null;
		role?: Role;
	};

	type Movie = {
		id: string;
		title: string;
		original_title?: string | null;
		thumbnail?: string | null;
		synopsis?: string | null;
		rating?: number | null;
		duration?: number | null;
		year?: number | null;
		date_published?: Date | null;
		label_quality?: string | null; // '1080p' | '720p' | '480p' | '360p'
		genre?: string | null;
		original_page_url?: string | null;
	};
	type MovieDetail = Movie & {
		votes?: number | null;
		release_date?: Date | null;
		updated_at?: Date | null;
		player_url?: PlayerUrl[] | null;
		trailer_url?: string | null;
		director?: MoviePerson[] | null;
		movie_star?: MoviePerson[] | null;
		genres?: MovieGenre[] | null;
		countries?: MovieCountry[] | null;
		similar_movies?: Movie[] | null;
	};

	type PlayerUrl = {
		url?: string | null; // '1080p' | '720p' | '480p' | '360p'
		type?: string | null;
	};
	type MoviePerson = {
		name?: string | null;
		page_url?: string | null;
	};
	type MovieGenre = {
		name?: string | null;
		page_url?: string | null;
	};
	type MovieCountry = {
		name?: string | null;
		page_url?: string | null;
	};

	type SeriesDetail = Movie & {
		season_name?: string | null;
		status?: string | null;
		votes?: number | null;
		season_list?: SeasonList[] | null;
		release_date?: Date | null;
		updated_at?: Date | null;
		director?: MoviePerson[] | null;
		movie_star?: MoviePerson[] | null;
		genres?: MovieGenre[] | null;
		countries?: MovieCountry[] | null;
		similar_movies?: Movie[] | null;
	};
	type SeasonList = {
		current_season?: string | null;
		total_season?: number | null;
		episode_list?: EpisodeList[] | null;
	};
	type EpisodeList = {
		episode_number?: number | null;
		episode_url?: string | null;
		player_url?: PlayerUrl[] | null;
		trailer_url?: string | null;
	};
	type SeriesEpisode = {
		episode_number?: number | null;
		episode_url?: string | null;
		player_url?: PlayerUrl[] | null;
		trailer_url?: string | null;
		download_url?: string | null;
	};

	type Anime = {
		id: string;
		title: string;
		title_japanese?: string | null;
		original_page_url?: string | null;
		thumbnail?: string | null;
		score?: string | null;
		producer?: string | null;
		type?: string | null;
		status?: string | null;
		total_episodes?: number | null;
		duration?: string | null;
		release_date?: string | null;
		released_day?: string | null;
		studio?: string | null;
		rating?: string | null;
		description?: string | null;
		created_at: string;
		updated_at: string;
		deleted_at?: Date | null;
		genres?: AnimeGenre[] | null;
		episodes?: AnimeEpisode[] | null;
	};
	type AnimeGenre = {
		name?: string | null;
		url?: string | null;
	};
	type AnimeEpisode = {
		id: string;
		title?: string | null;
		player_url?: string | null;
		page_url?: string | null;
		posted_by?: string | null;
		previous_episode_url?: string | null;
		next_episode_url?: string | null;
		see_all_episodes_url?: string | null;
		release_date?: string | null;
		release_time?: string | null;
		episode_number?: string | null;
		created_at: string;
		updated_at: string;
		deleted_at?: Date | null;
		list_episodes?: AnimeListOfEpisode[] | null;
		download_links?: DownloadLink[] | null;
	};
	type AnimeListOfEpisode = {
		name?: string | null;
		page_url?: string | null;
	};
	type DownloadLink = {
		name?: string | null;
		url?: string | null;
		size?: string | null;
		quality?: string | null; // '1080p' | '720p' | '480p' | '360p'
		format?: string | null;
	};

	type Setting = {
		id: string;
		scope?: string;
		key: string;
		value: string;
		description: string;
		group_name: string;
		created_at: string;
		updated_at: string;
	};
	type SettingsValue = {
		WEBSITE: SettingWeb;
		EMAIL: SettingEmail;
		SYSTEM: SettingSystem;
		MONETIZE: SettingMonetize;
	};
	type SettingWeb = {
		site_name?: string;
		site_tagline?: string;
		site_description?: string;
		site_keywords?: string;
		site_logo?: string;
		site_favicon?: string;
		site_email?: string;
		site_phone?: string;
		site_url?: string;
		site_created_at?: string;
	};
	type SettingEmail = {
		smtp_enabled?: boolean;
		smtp_service?: string;
		smtp_host?: string;
		smtp_port?: number;
		smtp_user?: string;
		smtp_password?: string;
		from_email?: string;
		from_name?: string;
	};
	type SettingSystem = {
		api_key?: string;
		theme?: "default" | "red-dark" | "green-dark" | "blue-dark";
		enable_documentation?: boolean;
		maintenance_mode?: boolean;
		maintenance_message?: string;
		source_logo_favicon: 'local' | 'remote';
		histats_tracking_code?: string;
		google_analytics_code?: string;
		play_store_app_url?: string;
		app_store_app_url?: string;
	};
	type SettingMonetize = {
		enable_monetize?: boolean;
		type_monetize?: 'adsense' | 'revenuecat' | 'adsterra';
		publisher_id?: string;
		enable_popup_ad?: boolean;
		enable_socialbar_ad?: boolean;
		auto_ad_code?: string;
		popup_ad_code?: string;
		socialbar_ad_code?: string;
		banner_rectangle_ad_code?: string;
		banner_horizontal_ad_code?: string;
		banner_vertical_ad_code?: string;
		native_ad_code?: string;
		direct_link_ad_code?: string;
	};

	type Application = {
		id: string;
		package_name?: string;
		key: string;
		value: string;
		description: string;
		group_name: string;
		created_at: string;
		updated_at: string;
	};
	type ApplicationValue = {
		CONFIG: ApplicationConfig;
		MONETIZE: ApplicationMonetize;
	};
	type ApplicationConfig = {
		name: string;
		api_key: string;
		package_name: string;
		version: string;
		type: string;
		store_url?: string | null;
		is_active: boolean;
	};
	type ApplicationMonetize = {
		enable_monetize: boolean;
		enable_admob: boolean;
		enable_unity_ad: boolean;
		enable_star_io_ad: boolean;
		enable_in_app_purchase: boolean;
		admob_id?: string | null;
		unity_ad_id?: string | null;
		star_io_ad_id?: string | null;
		admob_auto_ad?: string | null;
		admob_banner_ad?: string | null;
		admob_interstitial_ad?: string | null;
		admob_rewarded_ad?: string | null;
		admob_native_ad?: string | null;
		unity_banner_ad?: string | null;
		unity_interstitial_ad?: string | null;
		unity_rewarded_ad?: string | null;
		one_signal_id?: string | null;
	};

	type AnalyticsDaily = {
		date: string;
		count: number;
	};

	type Download = {
		id: string;
		user_id?: string;
		movie_id?: string;
		thumbnail_url?: string;
		title?: string;
		original_url?: string;
		status?: string;
		ip_address?: string;
		created_at?: string;
		updated_at?: string;
	};


	type ViewHistory = {
		id: string;
		user_id: string;
		name?: string | null;
		page_url?: string | null;
		ip_address?: string | null;
		user_agent?: string | null;
		browser_language?: string | null;
		device_type?: string | null;
		platform?: string | null;
		view_time?: Date | null;
		type?: 'movies' | 'series' | 'anime';
		created_at: string;
		updated_at: string;
	};

	type Ticket = {
		id: string;
		user_id: string;
		title?: string | null;
		description?: string | null;
		status?: 'open' | 'closed' | 'resolved';
		created_at: string;
		updated_at: string;
	};

	type Comment = {
		id: string;
		user_id: string;
		page_url: string;
		name?: string | null;
		email?: string | null;
		comment: string;
		type:
		| 'text'
		| 'image'
		| 'audio'
		| 'video'
		| 'document'
		| 'file'
		| 'location'
		| 'contact'
		| 'other';
		reply_to_id?: string | null;
		is_edited?: boolean;
		is_deleted?: boolean;
		media_url?: string | null;
		replies_count?: number;
		is_liked_by_me?: boolean;
		depth?: number;
		path?: string | null;
		created_at: string;
		updated_at: string;
		parent?: Comment | null;
		replies?: Comment[] | null;
		user?: User | null;
		likes?: LikeComment[] | null;
	};
	type LikeComment = {
		id: string;
		user_id: string;
		comment_id: string;
		created_at: string;
		updated_at: string;
		user?: User | null;
		comment?: Comment | null;
	};

	type Chat = {
		id: string;
		sender_id: string;
		receiver_id: string;
		message: string;
		type:
		| 'text'
		| 'image'
		| 'audio'
		| 'video'
		| 'document'
		| 'file'
		| 'location'
		| 'contact'
		| 'other';
		file_url?: string | null;
		file_name?: string | null;
		file_size?: number | null;
		mime_type?: string | null;
		is_read?: boolean;
		is_delivered?: boolean;
		is_edited?: boolean;
		is_deleted?: boolean;
		reply_to_id?: string | null;
		forwarded_from_id?: string | null;
		metadata?: any | null;
		created_at: string;
		updated_at: string;
		reply_to?: Chat | null;
		forwarded_from?: Chat | null;
		replies?: Chat[] | null;
		sender?: User | null;
		receiver?: User | null;
		reactions?: ChatReaction[] | null;
		reactions_count?: number;
		replies_count?: number;
		is_sent_by_me?: boolean;
		is_forwarded?: boolean;
		sender_name?: string | null;
		sender_avatar?: string | null;
		receiver_name?: string | null;
		receiver_avatar?: string | null;
	};
	type ChatReaction = {
		id: string;
		chat_id: string;
		user_id: string;
		reaction_type: 'like' | 'love' | 'haha' | 'wow' | 'sad' | 'angry';
		created_at: string;
		updated_at: string;
		user?: User | null;
		chat?: Chat | null;
	};
	type ChatMedia = {
		id: string;
		chat_id: string;
		url: string;
		type: 'image' | 'audio' | 'video' | 'document' | 'file' | 'location' | 'contact' | 'other';
		name?: string | null;
		size?: number | null;
		mime_type?: string | null;
		width?: number | null;
		height?: number | null;
		duration?: number | null;
		thumbnail?: string | null;
		created_at: string;
		chat?: Chat | null;
	};

	// ==========================================
	// Service Interfaces
	// ==========================================
	interface SettingService {
		registerSetting(
			settings: { key: string; scope: string; value: string; description?: string; group_name: string }[]
		): Promise<void | Error>;
		getPublicSettings(): Promise<SettingsValue | Error>;
		getAllSettings(): Promise<Setting[] | Error>;
		updateBulkSetting(
			settings: { key: string; value: string; description?: string; group_name: string }[]
		): Promise<void | Error>;
		updateFavicon(favicon: File): Promise<string | Error>;
		updateLogo(logo: File): Promise<string | Error>;
	}
	interface AuthService {
		Login(idToken: string): Promise<FirebaseAuthResponse | Error>;
		Register(data: RegisterSchema): Promise<FirebaseAuthResponse | Error>;
		ResetPassword(data: ResetPasswordSchema): Promise<string | Error>;
		VerifyEmail(data: VerifyEmailSchema): Promise<string | Error>;
		ChangePassword(data: ChangePasswordSchema): Promise<string | Error>;
	}
	interface AdminService {
		SetRole(data: SetRoleSchema): Promise<string | Error>;
		SearchUser(query: QueryParams): Promise<PaginatedResult<User>>
		HardDeleteUser(id: string): Promise<string | Error>;
		BulkDeleteUser(ids: string[]): Promise<string | Error>;
	}
	interface UserService {
		UpdateProfile(data: UpdateProfileSchema): Promise<User | Error>;
		CurrentUser(): Promise<User | Error>;
		UpdateAvatar(file: File): Promise<string | Error>;
		UpdatePassword(data: UpdatePasswordSchema): Promise<string | Error>;
		Logout(): Promise<string | Error>;
	}
	interface MovieService {
		GetHome(): Promise<HomeScrapperResponse | Error>;
		GetMoviesByGenre(slug: string, page: number): Promise<PaginatedResult<Movie>>;
		GetMoviesByCountry(country: string, page: number): Promise<PaginatedResult<Movie>>;
		GetMoviesByYear(year: number, page: number): Promise<PaginatedResult<Movie>>;
		SearchMovies(query: string, page: number): Promise<PaginatedResult<Movie>>;
		GetSpecialPage(path: string, page: number): Promise<PaginatedResult<Movie>>;
		GetMovieDetail(slug: string): Promise<MovieDetail | null>;
	}
	interface SeriesService {
		GetSeriesHome(): Promise<HomeScrapperResponse | Error>;
		GetSeriesByGenre(slug: string, page: number): Promise<PaginatedResult<Movie>>;
		GetSeriesByCountry(country: string, page: number): Promise<PaginatedResult<Movie>>;
		GetSeriesByYear(year: number, page?: number): Promise<PaginatedResult<Movie>>;
		SearchSeries(query: string, page?: number): Promise<PaginatedResult<Movie>>;
		GetSeriesByFeature(type: string, page?: number): Promise<PaginatedResult<Movie>>;
		GetSeriesSpecialPage(path: string, page?: number): Promise<PaginatedResult<Movie>>;
		GetSeriesDetail(slug: string): Promise<SeriesDetail | null>;
		GetSeriesEpisode(url: string): Promise<SeriesEpisode | null>;
	}
	interface AnimeService {
		GetLatest(page: number): Promise<PaginatedResult<Anime>>;
		Search(q: string, page: number): Promise<PaginatedResult<Anime>>;
		GetOngoing(page: number): Promise<PaginatedResult<Anime>>;
		GetGenres(): Promise<AnimeGenre[]>;
		GetDetail(url: string): Promise<Anime>;
		GetEpisode(url: string): Promise<AnimeEpisode[]>;
	}
	interface ServerStatusService {
		GetServerHealth(): Promise<ServerHealthResponse | null>;
		GetServerLogs(page: number, limit: number): Promise<PaginatedResult<ServerLogsResponse> | null>;
		ClearServerLogs(): Promise<void | Error>;
	}
	interface ApplicationService {
		RegisterApplication(
			apps: { package_name: string; key: string; value: string; description?: string; group_name: string }[]
		): Promise<boolean>;
		UpdateApplication(
			package_name: string,
			apps: { key: string; value: string; description?: string; group_name: string }[]
		): Promise<boolean>
		Search(query: QueryParams): Promise<PaginatedResult<ApplicationResponse>>;
		GetByPackageName(packageName: string): Promise<ApplicationResponse | null>;
		Delete(packageName: string): Promise<boolean>;
		BulkDelete(packageNames: string[]): Promise<boolean>;
	}

	interface ClientService {
		SendContact(data: ContactSchema): Promise<void | Error>;
	}


	// ==========================================
	// Analytics Interfaces
	// ==========================================
	interface DashboardStats {
		total_users: number;
		total_apps: number;
		total_platforms: number;
		total_downloads: number;
		total_subscriptions: number;
		total_transactions: number;
	}
	interface DashboardData {
		stats: DashboardStats;
		analytics: AnalyticsDaily[];
		recent_downloads: Download[];
	}
	interface DashboardResponse {
		data: DashboardData;
		pagination: ApiPagination;
	}
	// ==========================================
	// Server Status Interfaces
	// ==========================================
	interface ServerHealthResponse {
		database: string;
		redis: string;
		time: string;
	}
	interface ServerLogsResponse {
		level: string;
		message: string;
		timestamp: string;
		caller?: string;
		app?: string;
		env?: string;
		sql?: string;
		method?: string;
		path?: string;
		status?: number;
		count?: number;
		duration?: string;
		sql?: string;
		port?: number;
		args?: string[];
		command?: string;
		pipeline_size?: number;
		ip?: string;
		latency?: string;
		user_agent?: string;
		error?: string;
		args?: Record<string, any>[];
	}
}

export { };
