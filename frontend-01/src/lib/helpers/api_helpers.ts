import { PUBLIC_API_URL } from '$env/static/public';
import type { RequestEvent } from '@sveltejs/kit';
import { BaseHelper } from './base_helper';

export class ApiClientHandler extends BaseHelper implements ApiClient {
	private baseUrl: string;

	constructor(protected readonly event: RequestEvent) {
		super(event);
		this.baseUrl = PUBLIC_API_URL;
		// Ensure baseUrl doesn't end with slash to avoid double slashes
		if (this.baseUrl.endsWith('/')) {
			this.baseUrl = this.baseUrl.slice(0, -1);
		}
	}
	private getSecureHeaders(): Headers {
		const headers = new Headers();
		const url = new URL(this.event.request.url);

		const clientHost =
			this.event.request.headers.get('host') ||
			this.event.request.headers.get('x-forwarded-host') ||
			this.event.request.headers.get('x-real-host') ||
			url.host;

		const clientProto =
			this.event.request.headers.get('x-forwarded-proto') ||
			this.event.request.headers.get('x-forwarded-protocol') ||
			(url.protocol ? url.protocol.replace(':', '') : 'https');

		const clientOrigin = `${clientProto.startsWith('http') ? '' : clientProto + '://'}${clientHost}`;

		headers.set('Host', clientHost || '');
		headers.set('X-Forwarded-Host', clientHost || '');
		headers.set('X-Forwarded-Proto', clientProto || '');
		headers.set(
			'X-Forwarded-For',
			this.event.request.headers.get('x-forwarded-for') ||
			this.event.request.headers.get('x-real-ip') ||
			''
		);

		headers.set('X-Real-IP', this.event.request.headers.get('x-real-ip') || '');
		headers.set('Origin', clientOrigin);
		headers.set('Referer', this.event.request.headers.get('referer') || url.href);
		headers.set('X-Requested-With', 'XMLHttpRequest');
		headers.set('User-Agent', this.event.request.headers.get('user-agent') || '');

		headers.set('X-Content-Type-Options', 'nosniff');
		headers.set('X-Frame-Options', 'DENY');
		headers.set('X-XSS-Protection', '1; mode=block');

		if (!headers.get('Origin')) {
			console.warn('[ApiHandler] Origin header missing, using fallback:', url.origin);
			headers.set('Origin', url.origin);
		}
		return headers;
	}
	private async getCsrfToken(headers: Headers): Promise<{ token: string | null, cookie?: string }> {
		try {
			const cookieString = this.getCookieString();
			if (cookieString) {
				headers.set('Cookie', cookieString);
			}
			headers.set('X-Platform', 'browser');
			const response = await fetch(`${this.baseUrl}/token/csrf`, {
				method: 'GET',
				headers,
				credentials: 'include'
			});

			if (!response.ok) {
				throw new Error(`HTTP ${response.status}`);
			}

			let newCookie = '';
			const setCookie = response.headers.get('set-cookie');
			if (setCookie) {
				const csrfCookieMatch = setCookie.match(/csrf_session_id=([^;]+)/);
				if (csrfCookieMatch) {
					newCookie = `csrf_session_id=${csrfCookieMatch[1]}`;
				}
			}

			const data: ApiResponse<{ csrf_token: string }> = await response.json();
			return { token: data.data?.csrf_token ?? null, cookie: newCookie };
		} catch (err) {
			console.error('❌ Failed to get CSRF token', err);
			return { token: null, cookie: '' };
		}
	}
	private getCookieString(): string {
		const cookies: string[] = [];

		for (const [key, value] of Object.entries(this.event.cookies.getAll())) {
			cookies.push(`${key}=${value}`);
		}

		return cookies.join('; ');
	}
	private async createApiRequest<T>(
		method: HttpMethod,
		path: string,
		options: {
			data?: any;
			auth?: boolean;
			csrfProtected?: boolean;
			headers?: Record<string, string>;
		}
	): Promise<ApiResponse<T>> {
		// console.log('🔍 [ApiClient] Request:', method, path, options);
		const headers = this.getSecureHeaders();
		headers.set('X-Platform', 'browser');

		if (!(options.data instanceof FormData)) {
			headers.set('Content-Type', 'application/json');
		}

		if (options.headers) {
			for (const [key, value] of Object.entries(options.headers)) {
				headers.set(key, value);
			}
		}

		if (options.auth) {
			// Try getting token from options first (if passed explicitly), then cookie
			let accessToken = this.event.cookies.get('access_token');

			if (accessToken) {
				headers.set('Authorization', `Bearer ${accessToken}`);
			} else {
				console.warn('🔍 [ApiClient] Warning: Auth required but no access_token found in cookies');
			}
		}

		let cookieString = this.getCookieString();

		if (options.csrfProtected && method !== 'GET') {
			const { token, cookie } = await this.getCsrfToken(new Headers(headers));

			if (token) {
				headers.set('X-XSRF-TOKEN', token);
			}

			if (cookie) {
				if (cookieString) {
					// Remove old csrf_session_id to avoid conflicts
					cookieString = cookieString.replace(/csrf_session_id=[^;]+(; )?/, '');
					// Append new cookie
					cookieString = cookieString ? `${cookieString}; ${cookie}` : cookie;
				} else {
					cookieString = cookie;
				}
			}
		}

		if (cookieString) {
			headers.set('Cookie', cookieString);
		}

		const fullUrl = `${this.baseUrl}${path}`;

		let requestBody: any;
		if (method !== 'GET' && options.data) {
			requestBody = options.data instanceof FormData ? options.data : JSON.stringify(options.data);
		}

		try {
			const response = await fetch(fullUrl, {
				method,
				headers,
				body: requestBody,
				credentials: 'include'
			});

			const rawText = await response.text();
			let responseData: ApiResponse<T>;

			try {
				responseData = JSON.parse(rawText);
			} catch (e) {
				console.error('❌ Failed to parse JSON response. Raw text:', rawText);
				responseData = {
					status: response.status,
					success: false,
					message: `Invalid JSON response: ${response.status} ${response.statusText}`, // Include status text
					error: { code: 'INVALID_JSON', details: rawText } // Include raw text in error field for debugging
				};
			}

			const result: ApiResponse<T> = {
				status: responseData.status || response.status,
				success: responseData.success ?? response.ok,
				message: responseData.message || (response.ok ? 'Request successful' : 'Request failed'),
				data: responseData.data,
				pagination: responseData.pagination,
				error: responseData.error
			};

			if (!result.success) {
				console.error('❌ API Request returned error response', {
					url: fullUrl,
					method,
					status: result.status,
					message: result.message,
					error: result.error
				});
			}

			return result;
		} catch (error: any) {
			console.error('❌ API Request failed:', {
				url: fullUrl,
				method,
				error
			});
			return {
				status: 500,
				success: false,
				message: error.message || 'API request failed',
				error: {
					code: 'NETWORK_ERROR',
					details: process.env.NODE_ENV === 'development' ? error.stack : undefined
				}
			};
		}
	}
	private async createMultipartApiRequest<T>(
		method: HttpMethod,
		path: string,
		options: {
			data?: FormData;
			auth?: boolean;
			csrfProtected?: boolean;
			headers?: Record<string, string>;
		}
	): Promise<ApiResponse<T>> {
		const headers = this.getSecureHeaders();
		headers.set('X-Platform', 'browser');

		if (options.headers) {
			for (const [key, value] of Object.entries(options.headers)) {
				headers.set(key, value);
			}
		}

		if (options.auth) {
			let accessToken = this.event.cookies.get('access_token');
			if (accessToken) {
				headers.set('Authorization', `Bearer ${accessToken}`);
			}
		}

		let cookieString = this.getCookieString();

		if (options.csrfProtected && method !== 'GET') {
			const { token, cookie } = await this.getCsrfToken(new Headers(headers));

			if (token) {
				headers.set('X-XSRF-TOKEN', token);
			}

			if (cookie) {
				if (cookieString) {
					// Remove old csrf_session_id to avoid conflicts
					cookieString = cookieString.replace(/csrf_session_id=[^;]+(; )?/, '');
					// Append new cookie
					cookieString = cookieString ? `${cookieString}; ${cookie}` : cookie;
				} else {
					cookieString = cookie;
				}
			}
		}

		if (cookieString) {
			headers.set('Cookie', cookieString);
		}

		let requestBody: FormData = new FormData();
		if (method !== 'GET' && options.data) {
			requestBody = options.data;
		}

		try {
			const response = await fetch(`${this.baseUrl}${path}`, {
				method,
				headers,
				body: requestBody,
				credentials: 'include'
			});

			const responseData: ApiResponse<T> = await response.json().catch(() => ({
				status: response.status,
				success: false,
				message: 'Invalid JSON response'
			}));

			return {
				status: responseData.status || response.status,
				success: responseData.success ?? response.ok,
				message: responseData.message || (response.ok ? 'Request successful' : 'Request failed'),
				data: responseData.data,
				pagination: responseData.pagination,
				error: responseData.error
			};
		} catch (error: any) {
			console.error('❌ Multipart API Request failed:', error);
			return {
				status: 500,
				success: false,
				message: error.message || 'API request failed',
				error: {
					code: 'NETWORK_ERROR',
					details: process.env.NODE_ENV === 'development' ? error.stack : undefined
				}
			};
		}
	}
	public async authRequest<T>(
		method: HttpMethod,
		path: string,
		data?: any,
		headers?: Record<string, string>
	): Promise<ApiResponse<T>> {
		return this.createApiRequest<T>(method, path, {
			data,
			auth: true,
			csrfProtected: true,
			headers
		});
	}
	public async publicRequest<T>(
		method: HttpMethod,
		path: string,
		data?: any,
		csrfProtected?: boolean,
		headers?: Record<string, string>
	): Promise<ApiResponse<T>> {
		return this.createApiRequest<T>(method, path, {
			data,
			auth: false,
			csrfProtected: csrfProtected ?? false,
			headers
		});
	}
	public async multipartAuthRequest<T>(
		method: HttpMethod,
		path: string,
		data?: FormData,
		headers?: Record<string, string>
	): Promise<ApiResponse<T>> {
		return this.createMultipartApiRequest<T>(method, path, {
			data,
			auth: true,
			csrfProtected: true,
			headers
		});
	}
}
