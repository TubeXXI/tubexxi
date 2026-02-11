import { PostSchema } from '@/utils/schema.js';
import { SvelteComponent } from 'svelte';

declare global {
	interface Window {
		gc: NodeJS.GCFunction | undefined;
	}
	// component props
	interface PageMetaProps {
		path_url?: string;
		title?: string;
		tagline?: string;
		description?: string;
		keywords?: string[];
		robots?: string | boolean;
		canonical?: string;
		graph_type?: string;
		use_tagline?: boolean;
		language?: string;
		alternates?: {
			lang: string;
			href: string;
		}[];
	}
	interface CountryItem {
		name: string;
		code: string;
		emoji: string;
		unicode: string;
		image: string;
		dial_code: string;
		minLength: number;
		maxLength: number;
		regexPattern: string; // Regex pattern for phone number validation
	}
	interface TimezoneOption {
		zone: string;
		gmt: string;
		name: string;
	}
	interface ScrollAnimationConfig {
		threshold?: number;
		rootMargin?: string;
		animatedSelectors?: string[];
		autoInit?: boolean;
		isEnabled?: boolean;
	}
	interface ScrollAnimationState {
		isInitialized: boolean;
		isEnabled: boolean;
		observedElements: number;
		observer: IntersectionObserver | null;
	}
	type ToastMessage = {
		id: string;
		message: string;
		type: 'success' | 'error' | 'warning' | 'info';
		duration?: number;
	};
	interface MenuItem {
		id: number;
		title: string;
		url: string;
		icon?: string;
		child?: { title: string; url: string; icon?: string }[];
	}
	interface StatCard {
		title: string;
		description?: string;
		icon: string;
		value?: number;
		bgColor?: string;
		borderColor?: string;
		iconColor?: string;
		textColor?: string;
	}


	// api config
	interface ApiResponse<T = any, M extends Record<string, any> = ApiPagination> {
		status: number;
		success: boolean;
		message: string;
		data?: T | null;
		error?: ApiError;
		pagination?: M;
		headers?: Headers;
	}
	interface ApiPagination {
		current_page: number;
		limit: number;
		total_items: number;
		total_pages: number;
		has_prev: boolean;
		has_next: boolean;
		[key: string]: any;
	}
	interface ApiError {
		code: string;
		message?: string;
		redirect_url?: string;
		details?: any;
		retryable?: boolean;
		timestamp?: string;
		[key: string]: any;
	}
	interface ErrorResponse extends ApiResponse<undefined, undefined> {
		error: {
			code: string;
			details?: Record<string, unknown>;
			redirect_url?: Record<string, unknown>;
		};
	}
	type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE';
	interface FetchOptions {
		method?: HttpMethod;
		body?: any;
		headers?: HeadersInit;
		swrKey?: string;
		params?: Record<string, string | number | boolean | undefined>;
		timeout?: number;
		retries?: number;
		cache?: boolean;
	}
	interface ApiClient {
		authRequest<T>(
			method: HttpMethod,
			path: string,
			data?: any,
			headers?: Record<string, string>
		): Promise<ApiResponse<T>>;
		publicRequest<T>(
			method: HttpMethod,
			path: string,
			data?: any,
			csrfProtected?: boolean,
			headers?: Record<string, string>
		): Promise<ApiResponse<T>>;
		multipartAuthRequest<T>(
			method: HttpMethod,
			path: string,
			data?: FormData,
			headers?: Record<string, string>
		): Promise<ApiResponse<T>>;
	}

	// query params config
	interface QueryParamsConfig<T = any> {
		defaults: T;
		validators?: Partial<Record<keyof T, (value: any) => any>>;
	}
	interface SearchFieldConfig {
		field: string;
		type: 'string' | 'number' | 'date' | 'boolean';
		searchable?: boolean; // Default: true
	}
	const DEFAULT_PAGINATION: QueryParams = {
		page: 1,
		limit: 10,
		order_by: 'desc'
	};

	interface QueryParams {
		page: number;
		limit: number;
		search?: string;
		sort_by?: string;
		order_by?: 'asc' | 'desc';
		status?: string;
		include_deleted?: boolean;
		with_relations?: boolean;
		with_delete_column?: boolean;
		is_active?: boolean;
		is_verified?: boolean;
		user_id?: string;
		includes?: string[]; // e.g., ['user', 'category']
		fields?: string[]; // e.g., ['id', 'name', 'email']
		date_from?: Date | string;
		date_to?: Date | string;
		extra?: Record<string, any>;
		// for blog post
		tag?: string;
		series?: string;
		year?: number;
		month?: number;
		include_deleted?: boolean;
	}
	interface PaginatedResult<T> {
		data: T[];
		pagination: {
			current_page: number;
			total_pages: number;
			total_items: number;
			has_next: boolean;
			has_prev: boolean;
			limit: number;
			next_page?: number;
			prev_page?: number;
		};
	}
	interface QueryBuilderOptions<T> {
		searchFields?: string[]; // Fields to search in
		defaultSort?: string; // Default sort field
		customWhere?: any; // Custom where conditions
		defaultIncludes?: any; // Default includes
	}
	interface QueryState {
		page: number;
		limit: number;
		search: string;
		type: string;
		sort_by: string;
		order_by: 'asc' | 'desc';
		params: QueryParams;
		date_from: Date | string;
		date_to: Date | string;
	}
	interface CookieItem {
		lines: string[];
		path: string;
		valid: boolean;
	}
	interface BlogPost {
		meta: PostSchema
		component?: any;
		path: string
		readingTime: number
		words: number
		headings: {
			depth: number;
			value: string;
			slug: string;
		}[]
	}

}

export { };
